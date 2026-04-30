package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrRateLimited = errors.New("rate limited")

type UpstreamError struct {
	Code int
	Err  error
}

func (e *UpstreamError) Error() string { return fmt.Sprintf("upstream code=%d: %v", e.Code, e.Err) }
func (e *UpstreamError) Unwrap() error { return e.Err }

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// context error
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	// sentinel error
	if errors.Is(err, ErrRateLimited) {
		return true
	}

	// typed error
	var ue *UpstreamError
	if errors.As(err, &ue) {
		return ue.Code >= 500 && ue.Code <= 504
	}

	return false
}

func CallWithRetry(ctx context.Context) error {
	maxRetries := 3
	backoff := 100 * time.Millisecond // default waiting time

	for i := 0; i < maxRetries; i++ {
		err := callUpstream(ctx)
		if err == nil {
			return nil // 성공 시 즉시 반환
		}

		// 재시도 가능한 에러인지 확인
		if !IsRetryable(err) {
			return fmt.Errorf("non-retryable error: %w", err)
		}

		fmt.Printf("retry... (%d/%d): %v\n", i+1, maxRetries, err)

		// 다음 시도 전 대기
		select {
		case <-time.After(backoff):
			backoff *= 2 // 대기 시간을 늘려 서버 부담을 줄임
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return errors.New("max retries exceeded")
}

// error func
func callUpstream(ctx context.Context) error {
	return &UpstreamError{Code: 503, Err: errors.New("service unavailable")}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := CallWithRetry(ctx); err != nil {
		fmt.Printf("result error : %v \n", err)
	}
}
