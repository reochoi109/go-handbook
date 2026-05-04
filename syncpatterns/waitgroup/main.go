package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()

	jobs := []int{1, 2, 3, 4, 5}
	results := make(chan int)

	var wg sync.WaitGroup
	wg.Add(len(jobs))

	for _, job := range jobs {
		job := job
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			case <-time.After(80 * time.Millisecond):
			}
			select {
			case <-ctx.Done():
				return
			case results <- job * job:
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	sum := 0
	for v := range results {
		sum += v
	}
	fmt.Println("sum:", sum, "ctx:", ctx.Err())
}
