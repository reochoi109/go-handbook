package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(2)

	for i := 1; i <= 5; i++ {
		workerID := i
		g.Go(func() error {
			fmt.Printf("worker %d starting...\n", workerID)

			select {
			case <-time.After(1 * time.Second):
			case <-ctx.Done():
				fmt.Printf("worker %d canceled: %v\n", workerID, ctx.Err())
				return nil
			}

			// force an error
			if workerID == 3 {
				return fmt.Errorf("worker %d failed", workerID)
			}

			fmt.Printf("worker %d finishing\n", workerID)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Println("All tasks complete")
}
