package main

import (
	"fmt"
)

func main() {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for i := 1; i <= 5; i++ {
			ch <- i
		}
	}()

	sum := 0
	for v := range ch {
		sum += v
	}
	fmt.Printf("sum=%d\n", sum)

	buf := make(chan string, 2)
	buf <- "a"
	buf <- "b"
	close(buf)
	for s := range buf {
		fmt.Println("buf:", s)
	}
}
