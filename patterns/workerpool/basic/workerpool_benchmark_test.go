package main

import (
	"context"
	"sync"
	"testing"
)

func BenchmarkWorkerPool(b *testing.B) {
	for _, workers := range []int{1, 2, 4, 8} {
		b.Run("workers="+itoaBench(workers), func(b *testing.B) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			jobs := make(chan int, 1024)
			results := make(chan int, 1024)

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
							// simulate "some work" without sleeping
							results <- job * job
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
		})
	}
}

func BenchmarkBoundedVsUnbounded(b *testing.B) {
	b.Run("unbounded", func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		results := make(chan int, 1024)

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
			i := i
			go func() {
				defer wg.Done()
				select {
				case <-ctx.Done():
					return
				default:
				}
				results <- i * i
			}()
		}

		<-drainDone
	})

	for _, workers := range []int{4, 8} {
		b.Run("workerpool/workers="+itoaBench(workers), func(b *testing.B) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			jobs := make(chan int, 1024)
			results := make(chan int, 1024)

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
							results <- job * job
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
		})
	}
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
