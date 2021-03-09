package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

func fibonnaci(n uint64) uint64 {
	if n <= 1 || n > math.MaxInt64 {
		return n
	}
	return fibonnaci(n-1) + fibonnaci(n-2)
}

func double(n uint64) uint64 {
	return n * 2
}

func main() {
	var wg sync.WaitGroup

	key := []string{
		"emerald",
		"ruby",
		"sapphire",
	}

	/*
	 * Case 1
	 *
	 */
	time1 := time.Now()
	for _, v := range key {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			runWithMutexLock(
				func() {
					n := 0
					for i := 0; i < math.MaxInt32; i++ {
						n += i
					}
				},
			)
		}(v)
	}
	wg.Wait()
	duration := time.Since(time1)
	fmt.Printf("run with mutex lock time elapsed : %d\n", duration.Nanoseconds())

	time2 := time.Now()
	for _, v := range key {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			runWithMultiLock(
				key,
				func() {
					n := 0
					for i := 0; i < math.MaxInt32; i++ {
						n += i
					}
				},
			)
		}(v)
	}
	wg.Wait()
	duration2 := time.Since(time2)
	fmt.Printf("run with multi lock time elapsed : %d\n", duration2.Nanoseconds())

	/*
	 * Case 2
	 *
	 */
	n := 0
	m := 0
	o := 0

	key1 := "key1"
	key2 := "key2"
	key3 := "key3"

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			runWithMultiLock(
				key,
				func() {
					n += 2
				},
			)
		}(key1)
	}

	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			runWithMultiLock(
				key,
				func() {
					m += 2
				},
			)
		}(key2)
	}

	for k := 0; k < 10000; k++ {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			runWithMultiLock(
				key,
				func() {
					o += 2
				},
			)
		}(key3)
	}

	wg.Wait()

	fmt.Printf("n value : %d\n", n)
	fmt.Printf("m value : %d\n", m)
	fmt.Printf("o value : %d\n", o)
}
