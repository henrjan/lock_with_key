package main

import (
	"sync"
	"time"
)

var lock sync.Mutex

func runWithMutexLock(block func()) {
	lock.Lock()
	defer lock.Unlock()

	block()

	time.Sleep(100 * time.Millisecond)
}
