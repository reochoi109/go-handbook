package main

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"time"
)

const (
	addr            = "127.0.0.1:9090"
	maxPayloadSize  = 1 * 1024 * 1024 // 1MB로 제한 (보안)
	headerSize      = 4               // uint32 (4 bytes)
	networkDeadline = 5 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := runServer(ctx); err != nil && !errors.Is(err, net.ErrClosed) {
			fmt.Printf("[Server Err] %v\n", err)
		}
	}()

	time.Sleep(500 * time.Millisecond)
	if err := runClient(ctx); err != nil {
		fmt.Printf("[Client Err] %v\n", err)
	}
}

func runServer(ctx context.Context) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	fmt.Println("[Server] Listening on", addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go handleServerConn(conn)
	}
}

func handleServerConn(c net.Conn) {
	defer c.Close()
	fmt.Printf("[Server] Connected: %s\n", c.RemoteAddr())

	for {
		// read frame
		_ = c.SetReadDeadline(time.Now().Add(networkDeadline))
		msg, err := readFrame(c)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Printf("[Server] Read Error: %v\n", err)
			}
			return
		}

		fmt.Printf("[Server] Recv: %s\n", string(msg))

		_ = c.SetWriteDeadline(time.Now().Add(networkDeadline))
		resp := "ACK:" + string(msg)
		if err := writeFrame(c, []byte(resp)); err != nil {
			fmt.Printf("[Server] Write Error: %v\n", err)
			return
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

	payload := []byte("Hello Framing Protocol")

	// 전송
	_ = conn.SetWriteDeadline(time.Now().Add(networkDeadline))
	if err := writeFrame(conn, payload); err != nil {
		return err
	}

	// 응답 대기
	_ = conn.SetReadDeadline(time.Now().Add(networkDeadline))
	resp, err := readFrame(conn)
	if err != nil {
		return err
	}

	fmt.Printf("[Client] Recv: %s\n", string(resp))
	return nil
}

func writeFrame(w io.Writer, payload []byte) error {
	if len(payload) > maxPayloadSize {
		return fmt.Errorf("payload too large: %d", len(payload))
	}

	// 헤더와 페이로드를 담을 버퍼 한 번에 할당
	buf := make([]byte, headerSize+len(payload))
	binary.BigEndian.PutUint32(buf[:headerSize], uint32(len(payload)))
	copy(buf[headerSize:], payload)

	_, err := w.Write(buf)
	return err
}

func readFrame(r io.Reader) ([]byte, error) {
	// 1. header
	header := make([]byte, headerSize)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, err
	}

	// 2. length
	length := binary.BigEndian.Uint32(header)
	if length > maxPayloadSize {
		return nil, fmt.Errorf("frame size limit exceeded: %d", length)
	}

	// 3. payload
	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}

	return payload, nil
}
