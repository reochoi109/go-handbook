package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		deadline := time.Now().Add(250 * time.Millisecond)
		x := 0
		for time.Now().Before(deadline) {
			x = (x*33 + 7) % 1000003
		}
		_, _ = w.Write([]byte("ok\n"))
		_ = x
	})

	log.Println("pprof on http://localhost:6060/debug/pprof/")
	log.Println("work  on http://localhost:6060/work")
	log.Fatal(http.ListenAndServe("localhost:6060", nil))
}
