package node

import (
	"crypto/sha256"
	"crypto/tls"
	"net"
	"time"
)

type Connection struct{}

func (c *Connection) Send(b []byte) error {
	return nil
}

func (c *Connection) Receive() ([]byte, error) {
	return nil, nil
}

func (c *Connection) Close() error {
	return nil
}

type incomingPackage struct {
	b       []byte
	n       int
	offset  int
	rAddr   *net.UDPAddr
	next    *incomingPackage
	err     error
	cid     uint32
	pid     uint32
	prevPid uint32
}

type outgoingPackage struct {
	b    []byte
	n    int
	prev *outgoingPackage
	pKey [28]byte
	cid  uint32
}

type core struct {
	conn         *tls.Conn
	lastIncoming *incomingPackage
	incoming     *incomingPackage
	outgoing     *outgoingPackage
	heap         *heap
}

func (c *core) in(incoming *incomingPackage) bool {
	o := c.outgoing
	for o != nil {
		if incoming.cid == o.cid {
			copy(incoming.b[1432:], o.pKey[:])
			sig := sha256.Sum224(incoming.b[4:])
			if incoming.b[1] == sig[1] && incoming.b[2] == sig[2] && incoming.b[3] == sig[3] {
				goto CONTINUE
			}
		}
		o = o.prev
	}

	return false
CONTINUE:
	incoming.prevPid = bToID(incoming.b[7:10])
	incoming.pid = bToID(incoming.b[10:13])

	if c.lastIncoming.pid == incoming.prevPid {
		c.lastIncoming.next = incoming
		c.lastIncoming = incoming

		for incoming = c.heap.find(incoming.pid); incoming != nil; {
			c.lastIncoming.next = incoming
			c.lastIncoming = incoming
		}
	} else {
		c.heap.put(incoming)
	}

	return true
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

type heapItem struct {
	incoming *incomingPackage
	next     *heapItem
	prev     *heapItem
}

type heap struct {
	heap *heapItem
	last *heapItem
	len  int
	cap  int
}

func (h *heap) put(incoming *incomingPackage) {
	cur := h.heap
	for cur != nil {
		if cur.incoming.pid == incoming.pid {
			return
		}

		cur = cur.next
	}

	if h.cap <= h.len {
		if h.heap != nil {
			h.heap = h.heap.next
			if h.heap == nil {
				h.last = nil
			}
			h.len--
		}
	}

	item := &heapItem{
		incoming: incoming,
	}
	if h.last != nil {
		item.prev = h.last
		h.last.next = item
		h.last = item
	} else {
		h.heap = item
		h.last = item
	}
	h.len++
}

func (h *heap) find(pid uint32) *incomingPackage {
	cur := h.heap
	for cur != nil {
		if pid == cur.incoming.prevPid {
			if cur.prev != nil {
				cur.prev.next = cur.next
			} else {
				h.heap = cur.next
				if h.heap == nil {
					h.last = nil
				}
			}

			return cur.incoming
		}

		cur = cur.next
	}

	return nil
}

func bToID(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}
