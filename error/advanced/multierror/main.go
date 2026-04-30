package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"os"
)

func writeGzipFile(path string, data []byte) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	gz := gzip.NewWriter(f)

	defer func() {
		err = errors.Join(err, gz.Close(), f.Close())
	}()

	if _, err := gz.Write(data); err != nil {
		return fmt.Errorf("gzip write: %w", err)
	}

	return nil
}

func main() {
	fmt.Println(writeGzipFile("out.txt.gz", []byte("hello\n")))
}
