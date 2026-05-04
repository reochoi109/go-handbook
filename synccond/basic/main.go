package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var (
		mu    sync.Mutex
		cond  = sync.NewCond(&mu)
		ready bool
	)

	go func() {
		time.Sleep(120 * time.Millisecond)
		mu.Lock()
		ready = true
		mu.Unlock()

		// ready가 true가 됐다는 사실을 알림
		cond.Broadcast()
	}()

	mu.Lock()
	for !ready {
		cond.Wait()
	}
	fmt.Println("ready:", ready)
	mu.Unlock()
}
