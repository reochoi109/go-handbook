package main

import (
	"fmt"
	"sync"
	"time"
)

type NodeTask struct {
	ID    int
	Coord string
}

func processNode(task NodeTask) error {
	fmt.Printf("[worker] node %d saving... coord :%v\n", task.ID, task.Coord)
	time.Sleep(100 * time.Millisecond)
	return nil
}

func main() {
	jobs := make(chan NodeTask, 100)
	var wg sync.WaitGroup
	workers := 5
	for w := 0; w <= workers; w++ {
		go func(id int) {
			for task := range jobs {
				if err := processNode(task); err != nil {
					fmt.Printf("error : %v\n", err)
				}
				wg.Done()
			}
		}(w)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		jobs <- NodeTask{ID: i, Coord: "10.20"}
	}
	close(jobs)
	wg.Wait()
	fmt.Println("all node save")
}
