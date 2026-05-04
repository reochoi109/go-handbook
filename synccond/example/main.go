package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	q := NewQueue(2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 1; i <= 5; i++ {
			q.Push(i)
			fmt.Println("push", i)
			time.Sleep(30 * time.Millisecond)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 1; i <= 5; i++ {
			v := q.Pop()
			fmt.Println("pop ", v)
			time.Sleep(70 * time.Millisecond)
		}
	}()

	wg.Wait()
}

type Queue struct {
	mu       sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond

	buf  []int
	cap  int
	size int
	head int
	tail int
}

func NewQueue(capacity int) *Queue {
	q := &Queue{
		buf: make([]int, capacity),
		cap: capacity,
	}
	q.notEmpty = sync.NewCond(&q.mu)
	q.notFull = sync.NewCond(&q.mu)
	return q
}

func (q *Queue) Push(v int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for q.size == q.cap {
		q.notFull.Wait()
	}

	q.buf[q.tail] = v
	q.tail = (q.tail + 1) % q.cap
	q.size++

	q.notEmpty.Signal()
}

func (q *Queue) Pop() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	for q.size == 0 {
		q.notEmpty.Wait()
	}

	v := q.buf[q.head]
	q.head = (q.head + 1) % q.cap
	q.size--

	q.notFull.Signal()
	return v
}
