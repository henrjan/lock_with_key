package main

import (
	"sync"
	"testing"
)

var (
	div = 2
	c   = 10
)

func BenchmarkMultiLock(b *testing.B) {
	var wg sync.WaitGroup
	for n := 0; n < b.N; n++ {
		wg.Add(1)
		go func(key int) {
			runWithMultiLock(
				key,
				func() {
					fibonnaci(c)
				},
			)
			wg.Done()
		}(n)
	}
	wg.Wait()
}

func BenchmarkMutexLock(b *testing.B) {
	var wg sync.WaitGroup
	for n := 0; n < b.N; n++ {
		wg.Add(1)
		go func() {
			runWithMutexLock(
				func() {
					fibonnaci(c)
				},
			)
			wg.Done()
		}()
	}
	wg.Wait()
}
