package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	content := "User: Reo\nStatus: Backend Developer\nLog: Gzip compression test for technical blog.\n"
	filename := "data.gz"

	// 1. Compression Process
	if err := compressToFile(content, filename); err != nil {
		log.Fatal(err)
	}

	// 2. Decompression Process
	if err := decompressAndPrint(filename); err != nil {
		log.Fatal(err)
	}
}

func compressToFile(data string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Use gzip.BestCompression for maximum space saving
	zw, _ := gzip.NewWriterLevel(file, gzip.BestCompression)
	defer zw.Close()

	_, err = io.WriteString(zw, data)
	fmt.Printf("File saved: %s\n", filename)
	return err
}

func decompressAndPrint(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	zr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer zr.Close()

	if _, err := io.Copy(os.Stdout, zr); err != nil {
		return err
	}
	return nil
}
