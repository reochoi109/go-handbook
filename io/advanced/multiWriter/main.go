package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	inputData := "Backend developer Reo's Tech Log: io.MultiWriter Example Data"
	reader := strings.NewReader(inputData)

	consoleWriter := os.Stdout

	logFile, err := os.Create("process.log")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer logFile.Close()

	hashWriter := sha256.New()

	multi := io.MultiWriter(consoleWriter, logFile, hashWriter)

	written, err := io.Copy(multi, reader)
	if err != nil {
		log.Fatalf("Error occurred during writing: %v", err)
	}

	fmt.Printf("\nWriting completed (%d bytes)\n", written)

	hashInBytes := hashWriter.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)
	fmt.Printf("Generated SHA256 Hash: %s\n", hashString)
}
