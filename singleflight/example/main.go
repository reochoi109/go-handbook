package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

type Service struct {
	sf    singleflight.Group
	mu    sync.RWMutex
	cache map[string]string
	calls atomic.Int64
}

func NewService() *Service {
	return &Service{
		cache: make(map[string]string),
	}
}

func (s *Service) fetchFromDB(ctx context.Context, key string) (string, error) {
	s.calls.Add(1) // 실제 외부 호출 횟수를 기록합니다.

	select {
	case <-time.After(200 * time.Millisecond): // 의도 Latency 추가
		return fmt.Sprintf("DataFor:%s", key), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (s *Service) Get(ctx context.Context, key string) (string, error) {
	s.mu.RLock()
	val, ok := s.cache[key]
	s.mu.RUnlock()
	if ok {
		return val, nil
	}

	// 키 충돌을 방지
	vAny, err, _ := s.sf.Do("cachefill:"+key, func() (any, error) {
		s.mu.RLock()
		val, ok := s.cache[key]
		s.mu.RUnlock()
		if ok {
			return val, nil
		}

		data, err := s.fetchFromDB(ctx, key)
		if err != nil {
			return "", err
		}

		// 캐시
		s.mu.Lock()
		s.cache[key] = data
		s.mu.Unlock()

		return data, nil
	})

	if err != nil {
		return "", err
	}

	return vAny.(string), nil
}

func main() {
	service := NewService()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const concurrentRequests = 10
	key := "reo_777"

	var wg sync.WaitGroup
	wg.Add(concurrentRequests)

	fmt.Printf("--- Starting %d concurrent requests ---\n", concurrentRequests)
	startTime := time.Now()

	for i := 0; i < concurrentRequests; i++ {
		go func(id int) {
			defer wg.Done()
			val, err := service.Get(ctx, key)
			if err != nil {
				fmt.Printf("Goroutine %d error: %v\n", id, err)
				return
			}
			fmt.Printf("Goroutine %d result: %s\n", id, val)
		}(i)
	}

	wg.Wait()

	fmt.Println("---------------------------------")
	fmt.Printf("Total execution time: %v\n", time.Since(startTime))
	fmt.Printf("Actual DB fetch calls: %d\n", service.calls.Load())
	fmt.Println("---------------------------------")
}
