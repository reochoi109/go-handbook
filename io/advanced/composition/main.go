package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
)

func main() {
	combined := io.MultiReader(
		strings.NewReader("hello "),
		strings.NewReader("world"),
		strings.NewReader("\n"),
	)

	var tapped bytes.Buffer
	hash := sha256.New()

	// TeeReader로 읽히는 바이트를 side writer들에 복사
	tee := io.TeeReader(combined, io.MultiWriter(&tapped, hash))

	limited := io.LimitReader(tee, 11) // "hello world"만 읽기
	b, _ := io.ReadAll(limited)

	fmt.Print("read: " + string(b) + "\n")
	fmt.Print("tapped: " + tapped.String() + "\n")
	fmt.Printf("sha256: %x\n", hash.Sum(nil))
}
