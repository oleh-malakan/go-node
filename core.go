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

const (
	flags          = 0
	sigB1          = 1
	sigB2          = 2
	sigB3          = 3
	sigSumBegin    = 4
	cidBegin       = 4
	cidEnd         = 7
	pdidBegin      = 7
	pdidEnd        = 10
	didBegin       = 10
	didEnd         = 13
	datagramMinCap = 560
	datagramMaxCap = 1542
	datagramSigCap = 1460
)

type incomingDatagram struct {
	b       []byte
	n       int
	offset  int
	rAddr   *net.UDPAddr
	next    *incomingDatagram
	err     error
	cid     uint32
	did     uint32
	prevDid uint32
}

type outgoingDatagram struct {
	b    []byte
	n    int
	prev *outgoingDatagram
	pKey [sha256.Size224]byte
	cid  uint32
}

type container struct {
	conn *net.UDPConn
	in   func(*container, *incomingDatagram)

	next     *core
	inData   chan *incomingDatagram
	nextDrop chan *core
}

func (c *container) process() {
	go func() {
		for {
			i := &incomingDatagram{
				b: make([]byte, datagramSigCap),
			}
			i.n, i.rAddr, i.err = c.conn.ReadFromUDP(i.b)
			if i.err != nil {

				//continue
			}
			c.inData <- i
		}
	}()

	for {
		select {
		case i := <-c.inData:
			c.in(c, i)
		case d := <-c.nextDrop:
			c.next = d.next
			if c.next != nil {
				c.next.drop = c.nextDrop
				select {
				case c.next.reset <- nil:
				default:
				}
			}
		}
	}
}

type core struct {
	next     *core
	inData   chan *incomingDatagram
	nextDrop chan *core
	drop     chan *core
	reset    chan *struct{}
	isDrop   bool

	conn         *tls.Conn
	lastIncoming *incomingDatagram
	incoming     *incomingDatagram
	outgoing     *outgoingDatagram
	heap         *heap
}

func (c *core) process() {
	for !c.isDrop {
		select {
		case i := <-c.inData:
			if !c.in(i) && c.next != nil {
				c.next.inData <- i
			}
		case d := <-c.nextDrop:
			c.next = d.next
			if c.next != nil {
				c.next.drop = c.nextDrop
				select {
				case c.next.reset <- nil:
				default:
				}
			}
		}
	}

	for {
		select {
		case i := <-c.inData:
			if c.next != nil {
				c.next.inData <- i
			}
		case <-c.reset:
		case c.drop <- c:
			return
		}
	}
}

func (c *core) in(incoming *incomingDatagram) bool {
	o := c.outgoing
	for o != nil {
		if incoming.cid == o.cid {
			copy(incoming.b[incoming.n:sha256.Size224], o.pKey[:])
			sig := sha256.Sum224(incoming.b[sigSumBegin : incoming.n+sha256.Size224])
			if incoming.b[sigB1] == sig[sigB1] && incoming.b[sigB2] == sig[sigB2] && incoming.b[sigB3] == sig[sigB3] {
				goto CONTINUE
			}
		}
		o = o.prev
	}

	return false
CONTINUE:
	incoming.prevDid = bToID(incoming.b[pdidBegin:pdidEnd])
	incoming.did = bToID(incoming.b[didBegin:didEnd])

	if c.lastIncoming.did == incoming.prevDid {
		c.lastIncoming.next = incoming
		c.lastIncoming = incoming

		for incoming = c.heap.find(incoming.did); incoming != nil; {
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

const (
	heapCap = 256
)

type heapItem struct {
	incoming *incomingDatagram
	next     *heapItem
	prev     *heapItem
}

type heap struct {
	heap *heapItem
	last *heapItem
	len  int
}

func (h *heap) put(incoming *incomingDatagram) {
	cur := h.heap
	for cur != nil {
		if cur.incoming.did == incoming.did {
			return
		}

		cur = cur.next
	}

	if heapCap <= h.len {
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

func (h *heap) find(pid uint32) *incomingDatagram {
	cur := h.heap
	for cur != nil {
		if pid == cur.incoming.prevDid {
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
