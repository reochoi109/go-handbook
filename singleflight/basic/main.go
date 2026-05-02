package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

func main() {
	var g singleflight.Group

	var calls atomic.Int64
	work := func(ctx context.Context) (string, error) {
		calls.Add(1)
		select {
		case <-time.After(120 * time.Millisecond):
		case <-ctx.Done():
			return "", ctx.Err()
		}
		return "value-from-expensive-work", nil
	}

	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		i := i
		go func() {
			defer wg.Done()
			v, err, shared := g.Do("key:users:123", func() (any, error) {
				return work(ctx)
			})
			fmt.Printf("goroutine=%d value=%v err=%v shared=%v\n", i, v, err, shared)
		}()
	}
	wg.Wait()

	fmt.Println("expensive calls:", calls.Load())
}
