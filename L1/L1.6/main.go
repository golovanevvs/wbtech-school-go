package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	// через канал уведомления
	wg.Add(1)
	ch := make(chan struct{})
	go func() {
		defer wg.Done()
		byChannel(ch)
	}()
	time.Sleep(10 * time.Millisecond)
	close(ch)

	// через контекст при отмене
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer wg.Done()
		byCtxCancel(ctx)
	}()
	time.Sleep(10 * time.Millisecond)
	cancel()

	// через контекст по истечении заданного интервала времени
	wg.Add(1)
	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	go func() {
		defer wg.Done()
		byCtxTimeout(ctx)
	}()

	// через контекст с дедлайном
	wg.Add(1)
	deadLine := time.Now().Add(500 * time.Millisecond)
	ctx, cancel = context.WithDeadline(context.Background(), deadLine)
	defer cancel()
	go func() {
		defer wg.Done()
		byCtxDeadline(ctx)
	}()

	// через указанный промежуток времени с помощью time.After
	// возможна утечка памяти
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	byTimeAfter()
	// }()

	// через указанный промежуток времени по таймеру
	wg.Add(1)
	go func() {
		defer wg.Done()
		byTimer()
	}()

	// через специальную функцию runtime.Goexit()
	wg.Add(1)
	go func() {
		defer wg.Done()
		byRuntime()
	}()

	// через тикер
	wg.Add(1)
	go func() {
		defer wg.Done()
		byTicker()
	}()

	wg.Wait()
	fmt.Println("All goroutines finished")
}

func byChannel(ch chan struct{}) {
	for {
		select {
		case <-ch:
			fmt.Println("Stop by channel")
			return
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func byCtxCancel(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stop by context with cancel")
			return
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func byCtxTimeout(ctx context.Context) {
	<-ctx.Done()
	fmt.Println("Stop by context with timeout")
}

func byCtxDeadline(ctx context.Context) {
	<-ctx.Done()
	fmt.Println("Stop by context with deadline")
}

// нельзя остановить!
// func byTimeAfter() {
// 	select {
// 	case <-time.After(500 * time.Millisecond):
// 		fmt.Println("Stop by time.After")
// 		return
// 	}
// }

func byTimer() {
	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-timer.C:
		fmt.Println("Stop by timer")
		return
	default:
		time.Sleep(10 * time.Millisecond)
	}
}

func byRuntime() {
	defer fmt.Println("Stop by Goexit (defer)")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Stop by Goexit")
	runtime.Goexit()
}

func byTicker() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.After(500 * time.Millisecond)

	for {
		select {
		case t := <-ticker.C:
			fmt.Printf("Ticker tick at %v\n", t.Format("15:04:05.000"))
		case <-timeout:
			fmt.Println("Stop by ticker timeout")
			return
		}
	}
}
