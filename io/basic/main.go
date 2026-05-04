package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

func main() {
	// io 패키지는 데이터를 '어디서 읽어서 어디로 보낼지' 결정하는 도구 모음입니다.
	fmt.Println("-- 1.[copyDemo]")
	copyDemo()

	fmt.Println("-- 2.[teeDemo]")
	teeDemo()

	fmt.Println("-- 3.[limitDemo]")
	limitDemo()

	fmt.Println("-- 4.[bufioDemo]")
	bufioDemo()
}

func copyDemo() {

	src := strings.NewReader("Hello, Go World!\n") // (Reader)
	var dst bytes.Buffer                           // (Writer)

	// src에서 읽어 dst로 복사, 대용량 처리에 효율
	n, err := io.Copy(&dst, src)
	fmt.Printf("Copied bytes: %d, Error: %v\n", n, err)
	fmt.Print("Destination content: " + dst.String())
}

func teeDemo() {
	src := strings.NewReader("Important Data\n")
	var log bytes.Buffer

	r := io.TeeReader(src, &log) // src 읽을 때 마다 자동으로 log에 씀

	data, _ := io.ReadAll(r)
	fmt.Print("Read data: " + string(data))
	fmt.Print("Logged data: " + log.String())
}

func limitDemo() {
	src := strings.NewReader("1234567890")
	r := io.LimitReader(src, 4) // 최대 4바이트 : 예측 범위보다 큰 데이터 방지

	data, _ := io.ReadAll(r)
	fmt.Println("Limited read result:", string(data))
}

func bufioDemo() {
	var buf bytes.Buffer
	bw := bufio.NewWriter(&buf)
	bw.WriteString("Line A\n") // 입력 1
	bw.WriteString("Line B\n") // 입력 2

	bw.Flush()

	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		fmt.Println("Scanned line:", scanner.Text())
	}
}
