package main

import "time"

var mLock = NewMultipleLock()

func runWithMultiLock(key interface{}, block func()) {
	mLock.Lock(key)
	defer mLock.Unlock(key)

	block()

	time.Sleep(100 * time.Millisecond)
}
