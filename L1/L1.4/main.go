package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var n int
	var tSec int
	fmt.Println("Enter the number of workers")
	fmt.Scan(&n)
	fmt.Println("Enter the program runtime in seconds")
	fmt.Scan(&tSec)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(tSec)*time.Second)
	defer cancel()

	inCh := make(chan int)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT)

	wg := sync.WaitGroup{}

	for w := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for in := range inCh {
				fmt.Printf("worker %d received data: %d\n", i, in)
				time.Sleep(1 * time.Second)
			}
			fmt.Printf("worker %d has closed\n", i)
		}(w)
	}

	go func() {
		for range signalCh {
			cancel()
		}
	}()

	for i := 1; i > 0; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("received signal: %v\n", ctx.Err())
			close(inCh)
			wg.Wait()
			fmt.Println("end of work")
			return
		case inCh <- i:
		}
	}
}
