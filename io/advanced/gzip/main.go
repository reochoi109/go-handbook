package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

func main() {
	plain := []byte("hello gzip\nhello gzip\n")

	var compressed bytes.Buffer
	zw := gzip.NewWriter(&compressed)
	if _, err := zw.Write(plain); err != nil {
		panic(err)
	}
	// 중요: Close 해야 gzip 스트림이 완성됨
	if err := zw.Close(); err != nil {
		panic(err)
	}

	fmt.Println("compressed bytes:", compressed.Len())

	zr, err := gzip.NewReader(bytes.NewReader(compressed.Bytes()))
	if err != nil {
		panic(err)
	}
	defer zr.Close()

	out, err := io.ReadAll(zr)
	if err != nil {
		panic(err)
	}
	fmt.Print("decompressed:\n" + string(out))
}
