package memconn

import (
	"fmt"
	"net"
	"time"
)

// Package memconn은 UDP 예제에서 로컬 포트 바인딩 없이도
// net.PacketConn 스타일(ReadFrom/WriteTo)을 연습할 수 있게 해주는 in-memory 구현체입니다.

type Addr string

func (a Addr) Network() string { return "mem" }
func (a Addr) String() string  { return string(a) }

type packet struct {
	from net.Addr
	data []byte
}

type PacketConn struct {
	local  net.Addr
	peer   *PacketConn
	in     chan packet
	closed chan struct{}

	deadline time.Time
}

func NewPair() (*PacketConn, *PacketConn) {
	a := &PacketConn{local: Addr("A"), in: make(chan packet, 8), closed: make(chan struct{})}
	b := &PacketConn{local: Addr("B"), in: make(chan packet, 8), closed: make(chan struct{})}
	a.peer = b
	b.peer = a
	return a, b
}

func (c *PacketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	var timer <-chan time.Time
	if !c.deadline.IsZero() {
		d := time.Until(c.deadline)
		if d <= 0 {
			return 0, nil, fmt.Errorf("timeout")
		}
		timer = time.After(d)
	}

	select {
	case <-c.closed:
		return 0, nil, fmt.Errorf("closed")
	case <-timer:
		return 0, nil, fmt.Errorf("timeout")
	case pkt := <-c.in:
		n = copy(p, pkt.data)
		return n, pkt.from, nil
	}
}

func (c *PacketConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	if c.peer == nil {
		return 0, fmt.Errorf("no peer")
	}
	cp := make([]byte, len(p))
	copy(cp, p)
	select {
	case <-c.closed:
		return 0, fmt.Errorf("closed")
	case c.peer.in <- packet{from: c.local, data: cp}:
		return len(cp), nil
	}
}

func (c *PacketConn) Close() error {
	select {
	case <-c.closed:
		return nil
	default:
		close(c.closed)
		return nil
	}
}

func (c *PacketConn) LocalAddr() net.Addr { return c.local }

func (c *PacketConn) SetDeadline(t time.Time) error {
	c.deadline = t
	return nil
}

func (c *PacketConn) SetReadDeadline(t time.Time) error  { return c.SetDeadline(t) }
func (c *PacketConn) SetWriteDeadline(t time.Time) error { return c.SetDeadline(t) }
