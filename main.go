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
	count = 30
)

func main() {
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		operationWithMutex(count)
	}()

	go func() {
		defer wg.Done()
		operationWithMultiLock(count)
	}()

	wg.Wait()
}

func operationWithMutex(count int) {
	var wg sync.WaitGroup
	time1 := time.Now()
	for _, v := range key {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			runWithMutexLock(
				func() {
					for i := 0; i < count; i++ {
						fibonnaci(count)
					}
				},
			)
		}(v)
	}
	wg.Wait()
	duration := time.Since(time1)
	fmt.Printf("run with mutex lock time elapsed : %d\n", duration.Nanoseconds())
}

func operationWithMultiLock(count int) {
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
						fibonnaci(count)
					}
				},
			)
		}(v)
	}
	wg.Wait()
	duration2 := time.Since(time2)
	fmt.Printf("run with multi lock time elapsed : %d\n", duration2.Nanoseconds())
}

func fibonnaci(n int) uint64 {
	if n <= 1 || n > math.MaxInt64 {
		return uint64(n)
	}
	return fibonnaci(n-1) + fibonnaci(n-2)
}
