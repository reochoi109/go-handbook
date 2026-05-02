package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	lim := NewTokenBucket(5, 3) // rate=5/s, burst=3
	defer lim.Stop()

	for i := 1; i <= 12; i++ {
		if err := lim.Wait(ctx); err != nil {
			fmt.Println("stop:", err)
			return
		}
		fmt.Println("allow", i, "at", time.Now().Format("15:04:05.000"))
	}
}

type TokenBucket struct {
	tokens chan struct{}
	stop   chan struct{}
}

func NewTokenBucket(ratePerSec int, burst int) *TokenBucket {
	if ratePerSec <= 0 {
		ratePerSec = 1
	}
	if burst <= 0 {
		burst = 1
	}

	tb := &TokenBucket{
		tokens: make(chan struct{}, burst),
		stop:   make(chan struct{}),
	}

	// 초기 토큰은 burst만큼 채워둠(초기 burst 허용)
	for i := 0; i < burst; i++ {
		tb.tokens <- struct{}{}
	}

	interval := time.Second / time.Duration(ratePerSec)
	t := time.NewTicker(interval)

	go func() {
		defer t.Stop()
		for {
			select {
			case <-tb.stop:
				return
			case <-t.C:
				// 버킷이 가득 차면 refill은 drop
				select {
				case tb.tokens <- struct{}{}:
				default:
				}
			}
		}
	}()

	return tb
}

func (tb *TokenBucket) Stop() {
	select {
	case <-tb.stop:
	default:
		close(tb.stop)
	}
}

func (tb *TokenBucket) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-tb.tokens:
		return nil
	}
}
