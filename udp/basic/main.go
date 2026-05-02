package main

import (
	"fmt"
	"time"

	"github.com/reochoi109/go-handbook/udp/internal/memconn"
)

func main() {
	a, b := memconn.NewPair()

	go func() {
		// server role (b)
		buf := make([]byte, 256)
		n, addr, _ := b.ReadFrom(buf)
		fmt.Println("server recv from", addr.String(), "msg:", string(buf[:n]))

		_, _ = b.WriteTo([]byte("pong"), addr)
	}()

	// client role (a)
	_ = a.SetDeadline(time.Now().Add(500 * time.Millisecond))
	_ = b.SetDeadline(time.Now().Add(500 * time.Millisecond))

	_, _ = a.WriteTo([]byte("ping"), b.LocalAddr())

	buf := make([]byte, 256)
	n, addr, err := a.ReadFrom(buf)
	fmt.Println("client recv from", addr.String(), "msg:", string(buf[:n]), "err:", err)
}
