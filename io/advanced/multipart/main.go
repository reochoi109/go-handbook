package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
)

func main() {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	// field
	_ = w.WriteField("name", "reo")

	// file part (가짜)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="hello.txt"`)
	h.Set("Content-Type", "text/plain")
	part, _ := w.CreatePart(h)
	_, _ = part.Write([]byte("hello multipart\n"))

	_ = w.Close()

	fmt.Println("content-type:", w.FormDataContentType())
	fmt.Println("bytes:", body.Len())

	// parse
	r := multipart.NewReader(bytes.NewReader(body.Bytes()), w.Boundary())
	for {
		p, err := r.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		b, _ := io.ReadAll(p)
		fmt.Printf("part name=%q filename=%q data=%q\n", p.FormName(), p.FileName(), string(b))
		_ = p.Close()
	}
}
