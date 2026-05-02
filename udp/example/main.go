package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/reochoi109/go-handbook/udp/internal/memconn"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 700*time.Millisecond)
	defer cancel()

	client, server := memconn.NewPair()

	// server: request id를 그대로 돌려주는 echo
	go func() {
		buf := make([]byte, 256)
		for {
			n, from, err := server.ReadFrom(buf)
			if err != nil {
				return
			}
			// "ACK:" + payload
			_, _ = server.WriteTo(append([]byte("ACK:"), buf[:n]...), from)
		}
	}()

	// client: 타임아웃이면 재시도
	msg := "id=1 hello"
	resp, err := requestWithRetry(ctx, client, server.LocalAddr(), []byte(msg), 2, 120*time.Millisecond)
	fmt.Println("resp:", string(resp), "err:", err)
}

func requestWithRetry(
	ctx context.Context,
	c net.PacketConn,
	to net.Addr,
	payload []byte,
	maxAttempts int,
	timeout time.Duration,
) ([]byte, error) {
	if maxAttempts <= 0 {
		maxAttempts = 1
	}

	buf := make([]byte, 512)
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// send
		_, _ = c.WriteTo(payload, to)

		// recv with timeout (deadline)
		_ = c.SetReadDeadline(time.Now().Add(timeout))
		n, _, err := c.ReadFrom(buf)
		if err == nil {
			return buf[:n], nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}
	return nil, fmt.Errorf("no response after %d attempts", maxAttempts)
}
