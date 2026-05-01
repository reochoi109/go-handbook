package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	type job struct {
		id int
	}

	jobs := make(chan job)
	results := make(chan int)
	errs := make(chan error, 1)

	// producer
	go func() {
		defer close(jobs)
		for i := 1; i <= 10; i++ {
			jobs <- job{id: i}
		}
	}()

	// workers
	const workers = 4
	var wg sync.WaitGroup
	wg.Add(workers)

	var once sync.Once
	fail := func(err error) {
		if err == nil {
			return
		}
		once.Do(func() {
			errs <- err
			cancel()
		})
	}

	for w := 0; w < workers; w++ {
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case j, ok := <-jobs:
					if !ok {
						return
					}

					v, err := doWork(ctx, j)
					if err != nil {
						fail(err)
						return
					}

					select {
					case <-ctx.Done():
						return
					case results <- v:
					}
				}
			}
		}(w + 1)
	}

	// closer
	go func() {
		wg.Wait()
		close(results)
		close(errs)
	}()

	// consumer
	for v := range results {
		fmt.Println("result:", v)
	}

	select {
	case err := <-errs:
		if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("err:", err)
		}
	default:
	}
}

func doWork(ctx context.Context, j struct{ id int }) (int, error) {
	_ = ctx
	time.Sleep(70 * time.Millisecond)
	if j.id == 6 {
		return 0, fmt.Errorf("job failed id=%d", j.id)
	}
	return j.id * 10, nil
}
