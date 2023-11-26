package node

import (
	"crypto/tls"
	"net"
	"time"
)

type core struct {
	conn         *tls.Conn
	lastIncoming *incomingPackage
	incoming     *incomingPackage
	outgoing     *outgoingPackage
	heap         *heap
}

func (c *core) in(incoming *incomingPackage) {
	if compareID(c.lastIncoming.nextMac[0:32], incoming.b[33:65]) {
		c.lastIncoming.next = incoming
		c.lastIncoming = incoming
		var last *incomingPackage
		c.lastIncoming.next, last = c.heap.find(incoming.nextMac[0:32])
		if last != nil {
			c.lastIncoming = last
		}
	} else {
		c.heap.put(incoming)
	}
}

func (c *core) process() {
	select {
	default:
	}
}

func (c *core) Read(b []byte) (n int, err error) {
	return
}

func (c *core) Write(b []byte) (n int, err error) {
	return
}

func (c *core) Close() error {
	return nil
}

func (c *core) LocalAddr() net.Addr {
	return nil
}

func (c *core) RemoteAddr() net.Addr {
	return nil
}

func (c *core) SetDeadline(t time.Time) error {
	return nil
}

func (c *core) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *core) SetWriteDeadline(t time.Time) error {
	return nil
}
