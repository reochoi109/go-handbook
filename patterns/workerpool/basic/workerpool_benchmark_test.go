package main

import (
	"context"
	"sync"
	"testing"
)

func BenchmarkWorkerPool(b *testing.B) {
	for _, workers := range []int{1, 2, 4, 8} {
		b.Run("workers="+itoaBench(workers), func(b *testing.B) {
			runWorkerPool(b, workers, func(job int) int { return job * job })
		})
	}
}

func BenchmarkBoundedVsUnbounded(b *testing.B) {
	work := func(job int) int { return job * job }

	b.Run("unbounded", func(b *testing.B) { runUnbounded(b, work) })

	for _, workers := range []int{4, 8} {
		b.Run("workerpool/workers="+itoaBench(workers), func(b *testing.B) {
			runWorkerPool(b, workers, work)
		})
	}
}

func runWorkerPool(b *testing.B, workers int, work func(int) int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const bufSize = 1024
	jobs := make(chan int, bufSize)
	results := make(chan int, bufSize)

	var wg sync.WaitGroup
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobs:
					if !ok {
						return
					}
					results <- work(job)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	drainDone := make(chan struct{})
	go func() {
		for range results {
		}
		close(drainDone)
	}()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		jobs <- i
	}
	close(jobs)

	<-drainDone
}

func runUnbounded(b *testing.B, work func(int) int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const bufSize = 1024
	results := make(chan int, bufSize)

	var wg sync.WaitGroup
	wg.Add(b.N)

	go func() {
		wg.Wait()
		close(results)
	}()

	drainDone := make(chan struct{})
	go func() {
		for range results {
		}
		close(drainDone)
	}()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		job := i
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}
			results <- work(job)
		}()
	}

	<-drainDone
}

func itoaBench(n int) string {
	if n == 0 {
		return "0"
	}
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	var buf [32]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + (n % 10))
		n /= 10
	}
	return sign + string(buf[i:])
}
