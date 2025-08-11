package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type sample struct {
	data map[int]int
	mu   sync.Mutex
}

func main() {
	s := sample{
		data: make(map[int]int),
		mu:   sync.Mutex{},
	}

	wg := sync.WaitGroup{}

	for i := range 1000 {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			s.mu.Lock()
			s.data[v] = rand.Intn(1000)
			s.mu.Unlock()
		}(i)
	}

	wg.Wait()

	fmt.Println(s.data)
}
