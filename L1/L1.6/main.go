package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func main() {
	// через канал уведомления
	ch := make(chan struct{})
	go byChannel(ch)
	close(ch)

	// через контекст при отмене
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go byCtxCancel(ctx)
	cancel()

	// через контекст по истечении заданного интервала времени
	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	go byCtxTimeout(ctx)

	// через контекст по достижении заданного интервала времени
	deadLine := time.Now().Add(500 * time.Millisecond)
	ctx, cancel = context.WithDeadline(context.Background(), deadLine)
	defer cancel()
	go byCtxDeadline(ctx)

	// через указанный промежуток времени с помощью time.After
	go byTimeAfter()

	// через указанный промежуток времени по таймеру
	go byTimer()

	// через специальную функцию runtime.Goexit()
	go byRuntime()

	time.Sleep(5 * time.Second)
}

func byChannel(ch chan struct{}) {

	for {
		select {
		case <-ch:
			fmt.Println("Stop by channel")
			return
		}
	}
}

func byCtxCancel(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stop by context with cancel")
			return
		}
	}
}

func byCtxTimeout(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stop by context with timeout")
			return
		}
	}
}

func byCtxDeadline(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stop by context with deadline")
			return
		}
	}
}

func byTimeAfter() {
	timeToExit := time.After(500 * time.Millisecond)
	for {
		select {
		case <-timeToExit:
			fmt.Println("Stop by timeout of time.After")
			return
		}
	}
}
func byTimer() {
	t := time.NewTimer(500 * time.Millisecond)
	for {
		select {
		case <-t.C:
			fmt.Println("Stop by timeout of timer")
			return
		}
	}
}

func byRuntime() {
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Stop by Goexit")
	runtime.Goexit()
}
