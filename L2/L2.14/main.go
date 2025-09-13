package main

import (
	"fmt"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	orDone := make(chan interface{})
	if len(channels) == 0 {
		close(orDone)
		return orDone
	}

	go func() {
		defer close(orDone)
		switch len(channels) {
		case 1:
			<-channels[0]
		default:
			select {
			case <-channels[0]:
			case <-or(channels[1:]...):
			}
		}
	}()

	return orDone
}
func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}
