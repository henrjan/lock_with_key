package general

import (
	"sync"
)

type MultiLock struct {
	inUse sync.Map
	pool  *sync.Pool
	iLock sync.Mutex
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
	}
}

func (mLock *MultiLock) Lock(key interface{}) {
	lock := mLock.acquireLock(key)
	lock.Lock()
}

func (mLock *MultiLock) Unlock(key interface{}) {
	lock := mLock.releaseLock(key)
	lock.Unlock()
}

func (mLock *MultiLock) RLock(key interface{}) {
	lock := mLock.acquireLock(key)
	lock.RLock()
}

func (mLock *MultiLock) RUnlock(key interface{}) {
	lock := mLock.releaseLock(key)
	lock.RUnlock()
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

func (mLock *MultiLock) acquireLock(key interface{}) *sync.RWMutex {
	mLock.iLock.Lock()
	kl := mLock.getLock(key)
	kl.counter++
	mLock.iLock.Unlock()
	return kl.lock
}

func (mLock *MultiLock) releaseLock(key interface{}) *sync.RWMutex {
	mLock.iLock.Lock()
	kl := mLock.getLock(key)
	kl.counter--
	if kl.counter <= 0 {
		mLock.pool.Put(kl.lock)
		mLock.inUse.Delete(key)
	}
	mLock.iLock.Unlock()
	return kl.lock
}
