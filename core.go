package node

import (
	"crypto/sha256"
	"crypto/tls"
	"net"
	"time"
)

type core struct {
	conn         *tls.Conn
	lastIncoming *incomingPackage
	incoming     *incomingPackage
	iPKey        [32]byte
	outgoing     *outgoingPackage
	heap         *heap
}

func (c *core) in(incoming *incomingPackage) {
	if compare8(c.lastIncoming.b[25:33], incoming.b[17:25]) {
		for incoming != nil {
			incoming.b = append(incoming.b, c.iPKey[:]...)
			nPkey := sha256.Sum256(incoming.b[33:])
			if compare8(nPkey[:8], incoming.b[25:33]) {
				c.iPKey = nPkey
				c.lastIncoming.next = incoming
				c.lastIncoming = incoming
				incoming = c.heap.find(incoming.b[25:33])
			} else {
				incoming = c.heap.find(c.lastIncoming.b[25:33])
			}
		}
	} else {
		c.heap.put(incoming)
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
