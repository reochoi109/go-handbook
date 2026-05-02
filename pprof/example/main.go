package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	http.HandleFunc("/cpu", func(w http.ResponseWriter, r *http.Request) {
		// CPU burn (짧게)
		deadline := time.Now().Add(400 * time.Millisecond)
		n := 0
		for time.Now().Before(deadline) {
			n = (n*1103515245 + 12345) & 0x7fffffff
		}
		_, _ = fmt.Fprintf(w, "cpu ok n=%d\n", n)
	})

	http.HandleFunc("/alloc", func(w http.ResponseWriter, r *http.Request) {
		// 메모리 할당을 유발
		b := make([]byte, 2*1024*1024)
		for i := range b {
			b[i] = byte(rand.IntN(256))
		}
		_, _ = fmt.Fprintf(w, "alloc ok bytes=%d\n", len(b))
	})

	log.Println("server : http://localhost:6061/")
	log.Println("pprof  : http://localhost:6061/debug/pprof/")
	log.Println("try    : curl -s localhost:6061/cpu && curl -s localhost:6061/alloc")
	log.Fatal(http.ListenAndServe("localhost:6061", nil))
}
