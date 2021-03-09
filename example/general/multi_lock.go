package general

import (
	"sync"
)

type MultiLock struct {
	inUse   sync.Map
	pool    *sync.Pool
	acquire chan refKey
	release chan refKey
}

type keyLock struct {
	counter int
	lock    *sync.RWMutex
}

type refKey struct {
	key    interface{}
	lockCh chan *sync.RWMutex
}

func NewMultipleLock() *MultiLock {
	mLock := &MultiLock{
		pool: &sync.Pool{
			New: func() interface{} {
				return &sync.RWMutex{}
			},
		},
		acquire: make(chan refKey),
		release: make(chan refKey),
	}
	go mLock.lockUnlock()
	return mLock
}

func (mLock *MultiLock) Lock(key interface{}) {
	chGet := make(chan *sync.RWMutex)
	mLock.acquire <- refKey{key, chGet}
	lock := <-chGet
	lock.Lock()
}

func (mLock *MultiLock) Unlock(key interface{}) {
	chGet := make(chan *sync.RWMutex)
	mLock.release <- refKey{key, chGet}
	lock := <-chGet
	lock.Unlock()
}

func (mLock *MultiLock) RLock(key interface{}) {
	chGet := make(chan *sync.RWMutex)
	mLock.acquire <- refKey{key, chGet}
	lock := <-chGet
	lock.RLock()
}

func (mLock *MultiLock) RUnlock(key interface{}) {
	chGet := make(chan *sync.RWMutex)
	mLock.release <- refKey{key, chGet}
	lock := <-chGet
	lock.RUnlock()
}

func (mLock *MultiLock) lockUnlock() {
	for {
		select {
		case refKey := <-mLock.acquire:
			mLock.acquireLock(refKey.key, refKey.lockCh)
		case refKey := <-mLock.release:
			mLock.releaseLock(refKey.key, refKey.lockCh)
		default:
			// fmt.Println("Waiting...")
			// time.Sleep(time.Millisecond)
		}
	}
}

func (mLock *MultiLock) getLock(key interface{}) *keyLock {
	locker, _ := mLock.inUse.LoadOrStore(
		key,
		&keyLock{
			lock: mLock.pool.Get().(*sync.RWMutex),
		},
	)
	return locker.(*keyLock)
}

func (mLock *MultiLock) acquireLock(key interface{}, sendCh chan<- *sync.RWMutex) {
	kl := mLock.getLock(key)
	kl.counter++
	sendCh <- kl.lock
}

func (mLock *MultiLock) releaseLock(key interface{}, sendCh chan<- *sync.RWMutex) {
	kl := mLock.getLock(key)
	kl.counter--
	if kl.counter <= 0 {
		mLock.pool.Put(kl.lock)
		mLock.inUse.Delete(key)
	}
	sendCh <- kl.lock
}
