package main

import (
	"fmt"
	"math"
	"time"
)

type SlidingWindow struct {
	prevCount int // 이전 주기의 총 요청 수
	currCount int // 현재 주기의 총 요청 수
	limit     int
	lastReset time.Time     // 주기가 교체된 마지막 시점
	window    time.Duration // 주기
}

func NewSlidingWindow(limit int, window time.Duration) *SlidingWindow {
	return &SlidingWindow{
		limit:     limit,
		window:    window,
		lastReset: time.Now(),
	}
}

func (s *SlidingWindow) Allow() bool {
	now := time.Now()

	if now.Sub(s.lastReset) >= s.window {
		s.prevCount = s.currCount
		s.currCount = 0
		s.lastReset = now
	}

	// 2. 시간 흐름에 따른 가중치 계산
	passedFraction := float64(now.Sub(s.lastReset)) / float64(s.window)

	// 3. 현재 트래픽 추정치 계산
	estimate := float64(s.prevCount)*(1.0-passedFraction) + float64(s.currCount)

	// 4. 허용 여부 결정
	if int(math.Round(estimate)) < s.limit {
		s.currCount++
		return true
	}
	return false
}

func main() {
	limiter := NewSlidingWindow(10, time.Second)

	for i := 1; i <= 12; i++ {
		if limiter.Allow() {
			fmt.Printf("Request %d: ALLOWED (Current Count: %d)\n", i, limiter.currCount)
		} else {
			fmt.Printf("Request %d: REJECTED (Rate Limit Exceeded)\n", i)
		}
	}

	time.Sleep(1 * time.Second)        // wait sec
	limiter.Allow()                    // swap
	time.Sleep(250 * time.Millisecond) // wait 25%

	for i := 1; i <= 7; i++ {
		if limiter.Allow() {
			fmt.Printf("Request %d: ALLOWED (Current Count: %d)\n", i, limiter.currCount)
		} else {
			fmt.Printf("Request %d: REJECTED (Estimate reached limit)\n", i)
		}
	}
}
