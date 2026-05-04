package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	timeoutDemo()
	cancelDemo()
	nonBlockingDemo()
	tickerDemo()
}

func timeoutDemo() {
	ch := make(chan string)
	go func() {
		time.Sleep(150 * time.Millisecond)
		ch <- "done"
	}()

	select {
	case v := <-ch:
		fmt.Println("recv:", v)
	case <-time.After(80 * time.Millisecond):
		fmt.Println("timeout")
	}
}

func cancelDemo() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 1; i <= 5; i++ {
			time.Sleep(60 * time.Millisecond)
			ch <- i
		}
	}()

	// 2개만 받고 취소
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				fmt.Println("channel closed")
				return
			}
			fmt.Println("recv:", v)
			if v == 2 {
				cancel()
			}
		case <-ctx.Done():
			fmt.Println("canceled:", ctx.Err())
			return
		}
	}
}

func nonBlockingDemo() {
	ch := make(chan string)

	// default가 있으면 "아무 것도 준비되지 않았을 때" block하지 않습니다.
	select {
	case v := <-ch:
		fmt.Println("recv:", v)
	default:
		fmt.Println("no value (non-blocking)")
	}
}

func tickerDemo() {
	t := time.NewTicker(70 * time.Millisecond)
	defer t.Stop()

	deadline := time.After(230 * time.Millisecond)

	count := 0
	for {
		select {
		case <-t.C:
			count++
			fmt.Println("tick", count)
		case <-deadline:
			fmt.Println("stop")
			return
		}
	}
}
