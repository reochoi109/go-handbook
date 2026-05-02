package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := retry(ctx, 5, 30*time.Millisecond, func(ctx context.Context) error {
		return flakyCall(ctx)
	}); err != nil {
		fmt.Println("final err:", err)
	}

	fmt.Println()
	rateLimit(ctx, 4, func(i int) {
		fmt.Println("request", i, "at", time.Now().Format("15:04:05.000"))
	})
}

var errTemporary = errors.New("temporary")

func flakyCall(ctx context.Context) error {
	// 데모: 50% 확률로 실패
	select {
	case <-time.After(20 * time.Millisecond):
	case <-ctx.Done():
		return ctx.Err()
	}
	if rand.IntN(2) == 0 {
		return errTemporary
	}
	return nil
}

func retry(ctx context.Context, maxAttempts int, baseDelay time.Duration, fn func(context.Context) error) error {
	if maxAttempts <= 0 {
		maxAttempts = 1
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := fn(ctx); err != nil {
			lastErr = err
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return err
			}

			// 마지막 시도면 종료
			if attempt == maxAttempts {
				break
			}

			delay := backoffWithJitter(baseDelay, attempt)
			fmt.Println("attempt", attempt, "failed:", err, "-> sleep", delay)

			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
			continue
		}

		fmt.Println("attempt", attempt, "success")
		return nil
	}
	return lastErr
}

func backoffWithJitter(base time.Duration, attempt int) time.Duration {
	delay := base << (attempt - 1)
	max := 500 * time.Millisecond
	if delay > max {
		delay = max
	}
	j := 0.5 + rand.Float64()*0.5
	return time.Duration(float64(delay) * j)
}

func rateLimit(ctx context.Context, n int, fn func(i int)) {
	t := time.NewTicker(120 * time.Millisecond)
	defer t.Stop()

	for i := 1; i <= n; i++ {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			fn(i)
		}
	}
}
