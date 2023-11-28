package node

import (
	"crypto/sha256"
	"crypto/tls"
	"net"
	"time"
)

type incomingPackage struct {
	b      []byte
	n      int
	offset int
	rAddr  *net.UDPAddr
	next   *incomingPackage
	err    error
}

type outgoingPackage struct {
	b    []byte
	n    int
	prev *outgoingPackage
}

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

func compare8(a []byte, b []byte) bool {
	return a[0] == b[0] && a[1] == b[1] && a[2] == b[2] && a[3] == b[3] &&
		a[4] == b[4] && a[5] == b[5] && a[6] == b[6] && a[7] == b[7]
}

func compare16(a []byte, b []byte) bool {
	return a[0] == b[0] && a[1] == b[1] && a[2] == b[2] && a[3] == b[3] &&
		a[4] == b[4] && a[5] == b[5] && a[6] == b[6] && a[7] == b[7] &&
		a[8] == b[8] && a[9] == b[9] && a[10] == b[10] && a[11] == b[11] &&
		a[12] == b[12] && a[13] == b[13] && a[14] == b[14] && a[15] == b[15]
}
