package main

import (
	"fmt"
	"sync"
)

func main() {
	var (
		once sync.Once
		val  string
	)

	init := func() {
		val = "initialized"
		fmt.Println("init ran")
	}

	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			once.Do(init)
			fmt.Println("val:", val)
		}()
	}
	wg.Wait()
}
