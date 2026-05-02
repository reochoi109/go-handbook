package main

import (
	"fmt"
	"time"
)

func main() {
	limit := 3
	windowSize := time.Second

	var requests []time.Time

	// 요청 시뮬레이션
	testTimes := []int{0, 100, 200, 500, 1100, 1200, 1300} // 밀리초 단위

	for _, t := range testTimes {
		now := time.Now().Add(time.Duration(t) * time.Millisecond)

		boundary := now.Add(-windowSize)
		validIdx := 0
		for i, reqTime := range requests {
			if reqTime.After(boundary) {
				validIdx = i
				break
			}
			validIdx = i + 1
		}
		requests = requests[validIdx:]

		if len(requests) < limit {
			requests = append(requests, now)
			fmt.Printf("[%vms] ALLOW: Current window count: %d\n", t, len(requests))
		} else {
			fmt.Printf("[%vms] DENY : Rate limit exceeded\n", t)
		}
	}
}
