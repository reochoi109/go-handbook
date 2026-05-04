package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"time"
)

type event struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 파이프 생성 (In-memory 동기식 파이프)
	pr, pw := io.Pipe()

	var wg sync.WaitGroup
	wg.Add(2)

	// 1. Producer: 데이터를 생성하여 파이프에 쓰기
	go func() {
		defer wg.Done()

		var prodErr error
		defer func() {
			if prodErr != nil {
				pw.CloseWithError(prodErr)
			} else {
				pw.Close()
			}
		}()

		bw := bufio.NewWriterSize(pw, 1024) // 데이터 크기가 1024 이상 자동 pipe
		defer bw.Flush()                    // 수동

		enc := json.NewEncoder(bw)

		for i := 1; i <= 10; i++ {
			select {
			case <-ctx.Done():
				prodErr = ctx.Err()
				return
			default:
				msg := event{ID: i, Name: fmt.Sprintf("evt-%d", i)}
				if err := enc.Encode(msg); err != nil {
					prodErr = fmt.Errorf("encode fail: %w", err)
					return
				}

				time.Sleep(50 * time.Millisecond) // temp logic
			}
		}
	}()

	// 2. Consumer: 파이프에서 데이터를 읽어 처리하기
	go func() {
		defer wg.Done()
		defer pr.Close()

		dec := json.NewDecoder(pr)
		for {
			var e event
			// Decode는 내부적으로 필요한 만큼 Reader(pr)에서 데이터를 읽어옵니다.
			if err := dec.Decode(&e); err != nil {
				if err == io.EOF {
					log.Println("consumer: [SUCCESS] streaming finished")
					return
				}
				// Producer가 CloseWithError로 보낸 에러가 여기서 잡힙니다.
				log.Printf("consumer: [ERROR] %v\n", err)
				return
			}
			log.Printf("consumer: [RECEIVED] id=%d, name=%s\n", e.ID, e.Name)
		}
	}()

	wg.Wait()
	log.Println("main: all routines finished")
}
