package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(2)

	for i := 1; i <= 5; i++ {
		workerID := i
		g.Go(func() error {
			fmt.Printf("worker %d starting...\n", workerID)
			time.Sleep(1 * time.Second)

			fmt.Printf("worker %d finishing\n", workerID)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		fmt.Errorf("Error: %v\n", err)
	}
	fmt.Println("All tasks complete")
}
