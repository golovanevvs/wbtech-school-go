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

const numWorkers = 3

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inCh := make(chan int)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT)

	wg := sync.WaitGroup{}

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("worker %d has closed\n", i)
					return
				case in, ok := <-inCh:
					if !ok {
						fmt.Printf("worker %d has closed (channel closed)\n", i)
						return
					}
					fmt.Printf("worker %d received data: %d\n", i, in)
					time.Sleep(1 * time.Second)
				}
			}
		}(w)
	}

	go func() {
		<-signalCh
		fmt.Println("received SIGINT signal")
		cancel()
	}()

	i := 1
	for {
		select {
		case <-ctx.Done():
			close(inCh)
			wg.Wait()
			fmt.Println("end of work")
			return
		case inCh <- i:
			i++
		}
	}
}
