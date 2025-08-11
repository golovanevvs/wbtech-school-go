package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var n int
	fmt.Println("Enter the number of workers")
	fmt.Scan(&n)

	wg := sync.WaitGroup{}

	inCh := make(chan int)

	signalCh := make(chan os.Signal, 1)

	signal.Notify(signalCh, syscall.SIGINT)

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

	for i := 1; i > 0; i++ {
		select {
		case sig := <-signalCh:
			fmt.Printf("received signal: %v\n", sig)
			close(inCh)
			wg.Wait()
			fmt.Println("end of work")
			return
		default:
			inCh <- i
		}
	}

}
