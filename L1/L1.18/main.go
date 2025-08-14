package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type atomicExample struct {
	count int64
}

type mutexExample struct {
	count int64
	mu    sync.Mutex
}

func main() {
	numWorkers := 8
	n := 100

	atomicEx := atomicExample{}
	atomicEx.runCount(numWorkers, n)

	mutexEx := mutexExample{}
	mutexEx.runCount(numWorkers, n)
}

func (a *atomicExample) runCount(numWorkers int, n int) {
	printMu := sync.Mutex{}

	wg := sync.WaitGroup{}

	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func(i int) {
			defer wg.Done()
			for {
				cur := atomic.LoadInt64(&a.count)
				if cur >= int64(n) {
					break
				}
				if atomic.CompareAndSwapInt64(&a.count, cur, cur+1) {
					v := atomic.LoadInt64(&a.count)
					printMu.Lock()
					fmt.Printf("With atomic. Worker %d: %d\n", i, v)
					printMu.Unlock()
				}
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("With atomic. Count result: %d\n", a.count)
}

func (m *mutexExample) runCount(numWorkers int, n int) {
	wg := sync.WaitGroup{}

	wg.Add(numWorkers)

	for i := range numWorkers {
		go func(i int) {
			defer wg.Done()
			for {
				m.mu.Lock()
				if m.count < int64(n) {
					m.count++
					fmt.Printf("With mutex. Worker %d: %d\n", i, m.count)
					m.mu.Unlock()
				} else {
					m.mu.Unlock()
					break
				}
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("With mutex. Count result: %d\n", m.count)
}
