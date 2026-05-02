package main

import (
	"fmt"
	"time"
)

func main() {
	durationDemo()
	sinceUntilDemo()
	tickerDemo()
	timerReuseDemo()
}

// 1. 가독성
func durationDemo() {

	d := time.Minute
	fmt.Printf("set time: %v\n", d)
}

// 2. 시간 비교
func sinceUntilDemo() {

	start := time.Now()
	time.Sleep(30 * time.Millisecond)

	fmt.Printf("elapsed time (Since): %v\n", time.Since(start))

	deadline := time.Now().Add(50 * time.Millisecond)
	fmt.Printf("remaining time (Until): %v\n", time.Until(deadline))
}

// 3. 리소스 누수가 없는 Ticker 패턴
func tickerDemo() {

	ticker := time.NewTicker(40 * time.Millisecond)
	defer ticker.Stop()

	done := time.After(130 * time.Millisecond)

	count := 0
Loop:
	for {
		select {
		case t := <-ticker.C:
			count++
			fmt.Printf("tick %d at %v\n", count, t.Format("15:04:05.000"))
		case <-done:
			fmt.Println("ticker stopped safely")
			break Loop
		}
	}
}

// 4. GC 부담을 줄이는 Timer 재사용 패턴
func timerReuseDemo() {

	timer := time.NewTimer(35 * time.Millisecond)
	defer timer.Stop()

	for i := 1; i <= 3; i++ {
		t := <-timer.C
		fmt.Printf("Fired %d at %v\n", i, t.Format("15:04:05.000"))

		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(35 * time.Millisecond)
	}
	fmt.Println("timer reuse completed")
}
