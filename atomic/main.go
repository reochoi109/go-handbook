package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func checkStatus(isShutdown *atomic.Bool) {
	if isShutdown.Load() {
		fmt.Println("[System] Server is currently shut down.")
	} else {
		fmt.Println("[System] Server is running normally.")
	}
}

type Config struct{ Version string }

func updateConfig(config *atomic.Value, newVersion string) {
	config.Store(&Config{Version: newVersion})
}

func main() {
	var (
		count      int64
		isShutdown atomic.Bool
		config     atomic.Value
		wg         sync.WaitGroup
	)

	// init
	isShutdown.Store(false)
	config.Store(&Config{Version: "v1.0"})

	// counter
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&count, 1)
		}()
	}

	// status check
	checkStatus(&isShutdown)

	// config change
	updateConfig(&config, "v2.0")
	if cfg, ok := config.Load().(*Config); ok {
		fmt.Printf("[System] Config updated to version: %v\n", cfg.Version)
	}

	wg.Wait()
	fmt.Printf("[Result] Final count: %v\n", atomic.LoadInt64(&count))

}
