package main

import (
	"sync"
)

var lock sync.Mutex

func runWithMutexLock(block func()) {
	lock.Lock()
	defer lock.Unlock()

	block()
}
