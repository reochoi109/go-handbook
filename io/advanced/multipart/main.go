package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/textproto"
	"strings"
)

func main() {
	body, contentType, err := createMultipartData()
	if err != nil {
		log.Fatalf("Failed to generate multipart data: %v", err)
	}

	fmt.Printf("Content-Type: %s\n", contentType)
	fmt.Printf("Total Payload Size: %d bytes\n", body.Len())
	fmt.Println(strings.Repeat("-", 30))

	if err := parseMultipartData(body, contentType); err != nil {
		log.Fatalf("Failed to parse multipart data: %v", err)
	}
}

func createMultipartData() (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("username", "reo"); err != nil {
		return nil, "", err
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="attachment"; filename="reo_tech_log.txt"`)
	h.Set("Content-Type", "text/plain")

	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, "", err
	}

	content := "This is a multipart streaming test for Reo's Tech Log."
	if _, err := io.WriteString(part, content); err != nil {
		return nil, "", err
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}

func parseMultipartData(body *bytes.Buffer, contentType string) error {
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil || !strings.HasPrefix(mediaType, "multipart/") {
		return fmt.Errorf("invalid multipart content-type")
	}

	reader := multipart.NewReader(body, params["boundary"])

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		data, err := io.ReadAll(part)
		if err != nil {
			part.Close()
			return err
		}

		fmt.Printf("[Part Found]\n- FormName: %s\n- FileName: %s\n- Data: %s\n\n",
			part.FormName(), part.FileName(), string(data))
		part.Close()
	}
	return nil
}
