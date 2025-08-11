package main

import (
	"fmt"
	"time"
)

func main() {
	var tSec int
	fmt.Println("Enter the program runtime in seconds")
	fmt.Scan(&tSec)

	ch := make(chan int)

	go func() {
		for i := 1; i > 0; i++ {
			select {
			case _, ok := <-ch:
				if !ok {
					fmt.Println("Channel has closed")
					return
				}
			case ch <- i:
			}
		}
	}()

	go func() {
		<-time.After(time.Duration(tSec) * time.Second)
		close(ch)
	}()

	for v := range ch {
		fmt.Printf("Data: %d\n", v)
	}

	fmt.Println("End of work")
}
