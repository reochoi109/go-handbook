package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	srv, cli := net.Pipe()
	defer srv.Close()
	defer cli.Close()

	go func() {
		_ = server(srv)
	}()

	if err := client(cli); err != nil {
		fmt.Println("client err:", err)
	}
}

func server(c net.Conn) error {
	_ = c.SetDeadline(time.Now().Add(800 * time.Millisecond))

	msg, err := readFrame(c)
	if err != nil {
		return err
	}
	fmt.Println("server recv:", string(msg))

	// 응답도 프레임으로
	return writeFrame(c, []byte("ack:"+string(msg)))
}

func client(c net.Conn) error {
	_ = c.SetDeadline(time.Now().Add(800 * time.Millisecond))

	if err := writeFrame(c, []byte("hello")); err != nil {
		return err
	}

	resp, err := readFrame(c)
	if err != nil {
		return err
	}
	fmt.Println("client recv:", string(resp))
	return nil
}

func writeFrame(w io.Writer, payload []byte) error {
	// 프레임: [len(uint32)][payload...]
	var hdr [4]byte
	binary.BigEndian.PutUint32(hdr[:], uint32(len(payload)))

	bw := bufio.NewWriter(w)
	if _, err := bw.Write(hdr[:]); err != nil {
		return err
	}
	if _, err := bw.Write(payload); err != nil {
		return err
	}
	return bw.Flush()
}

func readFrame(r io.Reader) ([]byte, error) {
	br := bufio.NewReader(r)

	var hdr [4]byte
	if _, err := io.ReadFull(br, hdr[:]); err != nil {
		return nil, err
	}
	n := binary.BigEndian.Uint32(hdr[:])
	if n > 10*1024*1024 {
		return nil, fmt.Errorf("frame too large: %d", n)
	}

	buf := make([]byte, int(n))
	if _, err := io.ReadFull(br, buf); err != nil {
		return nil, err
	}
	return buf, nil
}
