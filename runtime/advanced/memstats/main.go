package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Println("1. before")
	printMem()

	s := make([][]byte, 0, 2000)
	for i := 0; i < 2000; i++ {
		s = append(s, make([]byte, 32*1024))
	}

	fmt.Println()
	fmt.Println("2. after alloc")
	printMem()

	// 메모리를 참조하지 않게 하고 GC 유도
	s = nil
	runtime.GC()

	time.Sleep(50 * time.Millisecond)

	fmt.Println()
	fmt.Println("3. after GC")
	printMem()
}

func printMem() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc=%dKB HeapAlloc=%dKB HeapObjects=%d NumGC=%d PauseTotal=%dms\n",
		m.Alloc/1024,
		m.HeapAlloc/1024,
		m.HeapObjects,
		m.NumGC,
		m.PauseTotalNs/1e6,
	)
}
