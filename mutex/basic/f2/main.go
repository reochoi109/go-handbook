package main

import (
	"fmt"
	"sync"
)

type SafeCounter struct {
	mu    sync.Mutex
	once  sync.Once
	count int
}

func NewSafeCounter() *SafeCounter {
	return &SafeCounter{}
}

func (c *SafeCounter) Init() {
	c.once.Do(func() {
		fmt.Println("init counter")
	})
}

func (c *SafeCounter) Inc() {
	c.Init()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *SafeCounter) GetCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func main() {
	counter := NewSafeCounter()
	jobs := make(chan int, 1000)
	var wg sync.WaitGroup
	workers := 10

	for w := 0; w < workers; w++ {
		go func() {
			for range jobs {
				counter.Inc()
				wg.Done()
			}
		}()
	}

	for j := 0; j < 1000; j++ {
		wg.Add(1)
		jobs <- j
	}

	close(jobs)
	wg.Wait()
	fmt.Printf("Final count : %d\n", counter.GetCount())
}
