package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

func main() {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
	defer cancel()

	// server side
	go func() {
		_ = handleServer(ctx, serverConn)
	}()

	// client side
	if err := runClient(ctx, clientConn); err != nil {
		fmt.Println("client err:", err)
	}
}

func handleServer(ctx context.Context, c net.Conn) error {
	_ = c.SetDeadline(time.Now().Add(600 * time.Millisecond))

	r := bufio.NewReader(c)
	line, err := r.ReadString('\n')
	if err != nil {
		return err
	}

	line = strings.TrimSpace(line)
	resp := "pong: " + line + "\n"

	_, err = c.Write([]byte(resp))
	return err
}

func runClient(ctx context.Context, c net.Conn) error {
	_ = c.SetDeadline(time.Now().Add(600 * time.Millisecond))
	if _, err := c.Write([]byte("ping\n")); err != nil {
		return err
	}

	r := bufio.NewReader(c)
	resp, err := r.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Print("resp: " + resp)
	_ = ctx
	return nil
}
