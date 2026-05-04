package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	var (
		mu   sync.Mutex
		cond = sync.NewCond(&mu)
		n    int
	)

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	// producer: n을 증가시키고 signal
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(90 * time.Millisecond)
			mu.Lock()
			n++
			mu.Unlock()
			cond.Signal()
		}
	}()

	// cancel watcher: ctx 취소되면 대기자 깨우기
	go func() {
		<-ctx.Done()
		cond.Broadcast()
	}()

	mu.Lock()
	defer mu.Unlock()

	for n < 5 && ctx.Err() == nil {
		cond.Wait()
	}

	fmt.Println("n:", n, "ctx:", ctx.Err())
}
