package general

import "sync"

type MultiLock struct {
	inUse sync.Map
	pool  *sync.Pool
	iLock sync.Map   // used for internal lock
	iPool *sync.Pool // used for internal lock
}

type keyLock struct {
	counter int
	lock    *sync.RWMutex
}

func NewMultipleLock() *MultiLock {
	return &MultiLock{
		pool: &sync.Pool{
			New: func() interface{} {
				return &sync.RWMutex{}
			},
		},
		iPool: &sync.Pool{
			New: func() interface{} {
				return &sync.Mutex{}
			},
		},
	}
}

func (mLock *MultiLock) Lock(key interface{}) {
	kl := mLock.getLock(key)
	kl.lock.Lock()
}

func (mLock *MultiLock) Unlock(key interface{}) {
	kl := mLock.releaseLock(key)
	kl.lock.Unlock()
}

func (mLock *MultiLock) RLock(key interface{}) {
	kl := mLock.getLock(key)
	kl.lock.RLock()
}

func (mLock *MultiLock) RUnlock(key interface{}) {
	kl := mLock.releaseLock(key)
	kl.lock.RUnlock()
}

func (mLock *MultiLock) getLock(key interface{}) *keyLock {
	iLock := mLock.useILock(key)
	defer mLock.releaseILock(key, iLock)

	locker, _ := mLock.inUse.LoadOrStore(
		key,
		&keyLock{
			lock: mLock.pool.Get().(*sync.RWMutex),
		},
	)
	kl := locker.(*keyLock)
	kl.counter++
	return kl
}

func (mLock *MultiLock) releaseLock(key interface{}) *keyLock {
	kl := mLock.getLock(key)

	kl.counter -= 2
	if kl.counter <= 0 {
		mLock.pool.Put(kl.lock)
		mLock.inUse.Delete(key)
	}
	return kl
}

func (mLock *MultiLock) useILock(key interface{}) *sync.Mutex {
	locker, _ := mLock.iLock.LoadOrStore(
		key,
		mLock.iPool.Get().(*sync.Mutex),
	)
	lock := locker.(*sync.Mutex)
	lock.Lock()
	return lock
}

func (mLock *MultiLock) releaseILock(key interface{}, lock *sync.Mutex) {
	lock.Unlock()
	mLock.iPool.Put(lock)
	mLock.iLock.Delete(key)
}
