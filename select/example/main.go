package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 700*time.Millisecond)
	defer cancel()

	a := producer(ctx, "A", 90*time.Millisecond)
	b := producer(ctx, "B", 140*time.Millisecond)

	readTwo(ctx, a, b)
	log.Println("done")
}

func producer(ctx context.Context, name string, every time.Duration) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		t := time.NewTicker(every)
		defer t.Stop()

		i := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				i++
				msg := fmt.Sprintf("%s-%d", name, i)
				select {
				case <-ctx.Done():
					return
				case out <- msg:
				}
			}
		}
	}()

	return out
}

func readTwo(ctx context.Context, a, b <-chan string) {
	for a != nil || b != nil {
		select {
		case <-ctx.Done():
			log.Println("canceled:", ctx.Err())
			return

		case v, ok := <-a:
			if !ok {
				a = nil // case 제거
				continue
			}
			log.Println("recv a:", v)

		case v, ok := <-b:
			if !ok {
				b = nil // case 제거
				continue
			}
			log.Println("recv b:", v)
		}
	}
}
