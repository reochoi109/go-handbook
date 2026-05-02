package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
	defer cancel()

	t := time.NewTicker(120 * time.Millisecond)
	defer t.Stop()

	for i := 1; i <= 10; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("stop:", ctx.Err())
			return
		case <-t.C:
			fmt.Println("allow request", i, "at", time.Now().Format("15:04:05.000"))
		}
	}
}
