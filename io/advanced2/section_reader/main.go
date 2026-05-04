package main

import (
	"bytes"
	"fmt"
	"io"
)

func main() {
	data := []byte("HEADER|payload:HELLO_WORLD|TRAILER")
	r := bytes.NewReader(data)

	off := int64(bytes.Index(data, []byte("HELLO_WORLD")))
	sec := io.NewSectionReader(r, off, int64(len("HELLO_WORLD")))

	b, _ := io.ReadAll(sec)
	fmt.Println(string(b))

	_, _ = sec.Seek(0, io.SeekStart)
	buf := make([]byte, 5)
	_, _ = sec.Read(buf)
	fmt.Println("first 5:", string(buf))
}
