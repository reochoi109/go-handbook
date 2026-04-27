package main

import (
	"flag"
	"fmt"
)

func main() {
	port := flag.Int("port", 8080, "server port")
	mode := flag.String("mode", "dev", "run mode")

	flag.Parse()

	fmt.Println("port:", *port)
	fmt.Println("mode:", *mode)
}
