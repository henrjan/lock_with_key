package main

var mLock = NewMultipleLock()

func runWithMultiLock(key interface{}, block func()) {
	mLock.Lock(key)
	defer mLock.Unlock(key)

	block()
}
