package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var m sync.Map

	// Store / Load
	m.Store("user:1", "Reo")
	if v, ok := m.Load("user:1"); ok {
		fmt.Println("load:", v)
	}

	// LoadOrStore
	v, loaded := m.LoadOrStore("user:1", "Other")
	fmt.Println("loadOrStore:", v, "loaded:", loaded)

	// Concurrent increment 예시(값을 atomic으로 저장)
	var hits atomic.Int64
	m.Store("hits", &hits)

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, _ := m.Load("hits")
			v.(*atomic.Int64).Add(1)
		}()
	}
	wg.Wait()

	v2, _ := m.Load("hits")
	fmt.Println("hits:", v2.(*atomic.Int64).Load())

	// Range: 순서 보장 없음
	m.Range(func(key, value any) bool {
		fmt.Printf("range: %v=%v\n", key, value)
		return true
	})
}
