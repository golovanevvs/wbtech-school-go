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

	timeout := time.After(time.Duration(tSec) * time.Second)

	go func() {
		i := 1
		for {
			select {
			case <-timeout:
				fmt.Println("Channel has closed")
				close(ch)
				return
			default:
				ch <- i
				i++
				time.Sleep(time.Second)
			}
		}
	}()

	for v := range ch {
		fmt.Printf("Data: %d\n", v)
	}

	fmt.Println("End of work")
}
