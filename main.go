package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

var (
	key = []string{
		"emerald",
		"ruby",
		"sapphire",
	}
	count = 40
)

func main() {
	var wg sync.WaitGroup

	/*
	 * running heavy operations using Mutex lock and multi lock
	 */

	wg.Add(2)
	go func() {
		defer wg.Done()
		operationWithMutex()
	}()

	go func() {
		defer wg.Done()
		operationWithMultiLock()
	}()

	wg.Wait()
}

func operationWithMutex() {
	var wg sync.WaitGroup
	time1 := time.Now()
	for _, v := range key {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			runWithMutexLock(
				func() {
					for i := 0; i < count; i++ {
						fibonnaci(uint64(i))
					}
				},
			)
		}(v)
	}
	wg.Wait()
	duration := time.Since(time1)
	fmt.Printf("run with mutex lock time elapsed : %d\n", duration.Milliseconds())
}

func operationWithMultiLock() {
	var wg sync.WaitGroup
	time2 := time.Now()
	for _, v := range key {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			runWithMultiLock(
				key,
				func() {
					for i := 0; i < count; i++ {
						fibonnaci(uint64(i))
					}
				},
			)
		}(v)
	}
	wg.Wait()
	duration2 := time.Since(time2)
	fmt.Printf("run with multi lock time elapsed : %d\n", duration2.Milliseconds())
}

func fibonnaci(n uint64) uint64 {
	if n <= 1 || n > math.MaxInt64 {
		return n
	}
	return fibonnaci(n-1) + fibonnaci(n-2)
}

func double(n uint64) uint64 {
	return n * 2
}
