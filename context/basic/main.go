package main

import (
	"context"
	"fmt"
	"time"
)

func doWork(ctx context.Context, name string) {
	fmt.Printf("[%s] process start\n", name)
	select {
	case <-time.After(2 * time.Second):
		fmt.Printf("[%s] process done\n", name)
	case <-ctx.Done():
		fmt.Printf("[%s] process stop %v\n", name, ctx.Err())
	}
}

func main() {
	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second)
	defer cancel1()
	doWork(ctx1, "Timeout Context")

	ctx2, cancel2 := context.WithCancel(context.Background())
	cancel2()
	doWork(ctx2, "Cancel Context")

	ctx3 := context.WithValue(context.Background(), "key", "value")
	doWork(ctx3, "Value Context")
}
