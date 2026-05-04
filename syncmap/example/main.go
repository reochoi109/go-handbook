package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	cache := &Cache{}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v := cache.GetOrCompute("k", func() string {
				time.Sleep(80 * time.Millisecond)
				return "value"
			})
			fmt.Println("goroutine", i, "got", v)
		}(i)
	}
	wg.Wait()
}

type Cache struct {
	m sync.Map
}

func (c *Cache) GetOrCompute(key string, compute func() string) string {
	if v, ok := c.m.Load(key); ok {
		return v.(string)
	}
	v := compute()
	c.m.Store(key, v)
	return v
}
