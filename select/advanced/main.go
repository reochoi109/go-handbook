package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	busyLoopFixed()
	timerReuse()
}

func busyLoopFixed() {
	fmt.Println("== busyLoopFixed ==")

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 1; i <= 3; i++ {
			time.Sleep(80 * time.Millisecond)
			ch <- i
		}
	}()

	// 안 좋은 예(설명):
	// for {
	//   select {
	//   case v := <-ch:
	//     ...
	//   default:
	//     // 아무 것도 안 하고 계속 루프 → busy loop
	//   }
	// }

	// 권장 패턴:
	// - default가 필요하다면 "sleep/backoff" 같은 완충이 있어야 함
	// - 혹은 tick/timeout을 case로 넣어 block하는 구조로 설계합니다.
	idle := time.NewTicker(30 * time.Millisecond)
	defer idle.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("stop:", ctx.Err())
			return
		case v, ok := <-ch:
			if !ok {
				fmt.Println("channel closed")
				return
			}
			fmt.Println("recv:", v)
		case <-idle.C:
			// 대기
			fmt.Println("idle tick")
		}
	}
}

func timerReuse() {
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	// time.After를 루프에서 계속 만들면, 매 반복마다 새로운 타이머가 생성
	// 짧은 루프/고빈도에서는 비용이 될 수 있습니다.
	//
	// 권장: time.NewTimer + Reset으로 재사용
	timer := time.NewTimer(60 * time.Millisecond)
	defer timer.Stop()

	i := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("stop:", ctx.Err())
			return
		case <-timer.C:
			i++
			fmt.Println("timer fired", i)

			// Reset 전에는 timer가 이미 만료되어 C가 비어 있는 상태여야 안전
			// 여기서는 <-timer.C로 소비했기 때문에 OK
			timer.Reset(60 * time.Millisecond)
		}
	}
}
