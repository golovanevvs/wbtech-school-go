package main

import (
	"fmt"
	"sync"
)

func main() {
	arr := [5]int{2, 4, 6, 8, 10}

	wg := sync.WaitGroup{}

	for _, v := range arr {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Printf("%d в квадрате равно %d\n", i, i*i)
		}(v)
	}

	wg.Wait()
}
