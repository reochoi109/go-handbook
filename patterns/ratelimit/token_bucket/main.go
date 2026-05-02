package main

import (
	"context"
	"fmt"
	"time"
)

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

	for i := 0; i < burst; i++ {
		tb.tokens <- struct{}{}
	}

	go tb.refillLoop(ratePerSec)

	return tb
}

func (tb *TokenBucket) refillLoop(ratePerSec int) {
	ticker := time.NewTicker(time.Second / time.Duration(ratePerSec))
	defer ticker.Stop()

	for {
		select {
		case <-tb.stop:
			return
		case <-ticker.C:
			select {
			case tb.tokens <- struct{}{}:
			default:
			}
		}
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

func (tb *TokenBucket) Allow() bool {
	select {
	case <-tb.tokens:
		return true
	default:
		return false // default 있기 때문에 계속 기다리지 않음
	}
}

func (tb *TokenBucket) Stop() {
	select {
	case <-tb.stop:
		return // 이미 닫힘
	default:
		close(tb.stop)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	lim := NewTokenBucket(5, 3) // 초당 5개 리필, 최대 3개 보관
	defer lim.Stop()

	for i := 1; i <= 10; i++ {
		if err := lim.Wait(ctx); err != nil {
			fmt.Println("Stop:", err)
			break
		}
		fmt.Printf("[%d] Allow at %s\n", i, time.Now().Format("15:04:05.000"))
	}
}
