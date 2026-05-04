package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"strings"
)

func main() {
	src := strings.NewReader(strings.Repeat("hello-", 10000))

	h := sha256.New()
	tee := io.TeeReader(src, hashWriter{h: h})

	n, err := io.Copy(io.Discard, tee)
	if err != nil {
		fmt.Println("copy err:", err)
		return
	}

	fmt.Println("bytes:", n)
	fmt.Println("sha256:", hex.EncodeToString(h.Sum(nil)))
}

type hashWriter struct {
	h hash.Hash
}

func (w hashWriter) Write(p []byte) (int, error) { return w.h.Write(p) }
