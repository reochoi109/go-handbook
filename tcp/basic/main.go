package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

const (
	addr         = "127.0.0.1:8080"
	readTimeout  = 2 * time.Second
	writeTimeout = 2 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runServer(ctx); err != nil && !errors.Is(err, net.ErrClosed) {
			fmt.Printf("[Server Err] %v\n", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runClient(ctx); err != nil {
			fmt.Printf("[Client Err] %v\n", err)
		}
	}()

	wg.Wait()
}

func runServer(ctx context.Context) error {
	lc := net.ListenConfig{}
	ln, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Printf("[Server] %s waiting...\n", addr)

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	remoteAddr := c.RemoteAddr().String()
	r := bufio.NewReader(c)

	const idleTimeout = 10 * time.Second

	for {
		_ = c.SetReadDeadline(time.Now().Add(idleTimeout))

		line, err := r.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Printf("[Server] client close : %s\n", remoteAddr)
			} else {
				fmt.Printf("[Server] error or timeout (%s): %v\n", remoteAddr, err)
			}
			return
		}

		cmd := strings.TrimSpace(strings.ToLower(line))
		if cmd == "ping" {
			fmt.Printf("[%s] %s -> RECEIVED: PING\n", time.Now().Format("15:04:05"), remoteAddr)

			_ = c.SetWriteDeadline(time.Now().Add(writeTimeout))
			_, _ = c.Write([]byte("PONG\n"))
		}
	}
}

func runClient(ctx context.Context) error {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	r := bufio.NewReader(conn)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("[Client] conntection close.")
			return nil
		case t := <-ticker.C:
			// 1. Ping 전송
			_ = conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			message := "ping\n"
			if _, err := conn.Write([]byte(message)); err != nil {
				return fmt.Errorf("write error: %w", err)
			}

			// 2. Pong 응답 대기
			_ = conn.SetReadDeadline(time.Now().Add(readTimeout))
			resp, err := r.ReadString('\n')
			if err != nil {
				return fmt.Errorf("read error: %w", err)
			}

			fmt.Printf("[Client] %s -> srv response: %s", t.Format("15:04:05"), resp)
		}
	}
}
