package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	server := &http.Server{Addr: ":8080"}

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("worker: working")
		<-ctx.Done()
		fmt.Println("worker: complete")
	}()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("server error: %v\n", err)
		}
	}()

	<-ctx.Done()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(timeoutCtx)
	wg.Wait()

	fmt.Println("system close")
}
