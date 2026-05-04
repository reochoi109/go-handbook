package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	header := strings.NewReader("--- START OF LOG ---\n")
	body := strings.NewReader("Task: Process Data\nStatus: Success\n")
	footer := strings.NewReader("--- END OF LOG ---\n")

	combinedReader := io.MultiReader(header, body, footer)

	fmt.Println("Reading from combined stream:")

	written, err := io.Copy(os.Stdout, combinedReader)
	if err != nil {
		log.Fatalf("Error occurred during reading: %v", err)
	}
	fmt.Printf("\nTotal bytes read: %d\n", written)
}
