package node

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"net"
	"time"
)

const (
	dataBegin      = 16
	datagramMinLen = 560
	datagramMaxLen = 1432
	datagramSigLen = 1460
)

type incomingDatagram struct {
	b       [datagramSigLen]byte
	dataEnd int
	offset  int
	rAddr   *net.UDPAddr
	next    *incomingDatagram
	err     error
	cid     uint32
	did     uint32
	prevDid uint32
}

type outgoingDatagram struct {
	b       [datagramSigLen]byte
	dataEnd int
	prev    *outgoingDatagram
	pKey    [sha256.Size224]byte
	cid     uint32
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
		var n int
		for {
			i := &incomingDatagram{}
			n, i.rAddr, i.err = c.conn.ReadFromUDP(i.b[:])
			if i.err != nil {

				//continue
			}
			i.dataEnd = n - 1
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
				case c.next.signal <- nil:
				default:
				}
			}
		}
	}
}

type core struct {
	isProcess bool
	next      *core
	nextDrop  chan *core
	drop      chan *core
	signal    chan *struct{}

	inData       chan *incomingDatagram
	lastIncoming *incomingDatagram
	incoming     *incomingDatagram
	outgoing     *outgoingDatagram
	heap         *heap

	conn       *tls.Conn
	tlsRead    tlsRead
	tlsProcess *tlsProcess
}

func (c *core) process() {
	ctx, cancelHandshake := context.WithCancel(context.Background())
	go c.conn.HandshakeContext(ctx)

	c.tlsProcess = &tlsProcess{
		tlsInData:   make(chan *incomingDatagram),
		tlsInSignal: make(chan *struct{}),
	}
	<-c.signal
	c.tlsRead = c.tlsProcess
	c.signal <- nil
	timerCancelHandshake := time.NewTimer(time.Duration(200) * time.Millisecond)
	for c.isProcess {
		select {
		case i := <-c.inData:
			if !c.in(i) && c.next != nil {
				c.next.inData <- i

				continue
			}
			c.tslIn()
		case <-c.tlsProcess.tlsInSignal:
			c.tslIn()
		case d := <-c.nextDrop:
			c.next = d.next
			if c.next != nil {
				c.next.drop = c.nextDrop
				select {
				case c.next.signal <- nil:
				default:
				}
			}
		case <-timerCancelHandshake.C:
			c.isProcess = false
		}
	}

	cancelHandshake()
	for {
		select {
		case i := <-c.inData:
			if c.next != nil {
				c.next.inData <- i
			}
		case <-c.signal:
		case c.drop <- c:
			return
		}
	}
}

func (c *core) in(incoming *incomingDatagram) bool {
	o := c.outgoing
	for o != nil {
		if incoming.cid == o.cid {
			sig := signing(incoming.b, incoming.dataEnd, o.pKey[:])
			if incoming.b[1] == sig[1] && incoming.b[2] == sig[2] &&
				incoming.b[3] == sig[3] && incoming.b[4] == sig[4] {
				goto CONTINUE
			}
		}
		o = o.prev
	}

	return false
CONTINUE:
	incoming.prevDid = prevDidFromB(incoming.b)
	incoming.did = didFromB(incoming.b)

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

type tlsRead interface {
	read(b []byte) (n int, err error)
}

type tlsProcess struct {
	tlsInData   chan *incomingDatagram
	tlsInSignal chan *struct{}
}

func (c *tlsProcess) read(b []byte) (n int, err error) {
	cur := <-c.tlsInData
	for cur.next != nil {
		//
		cur = cur.next
	}

	c.tlsInSignal <- nil

	return 0, nil
}

func (c *core) tslIn() {
	select {
	case c.tlsProcess.tlsInData <- c.incoming:
		c.incoming = c.lastIncoming
	default:
	}
}

func (c *core) Read(b []byte) (n int, err error) {
	return c.tlsRead.read(b)
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

type Stream struct{}

func (c *Stream) Send(b []byte) error {
	return nil
}

func (c *Stream) Receive() ([]byte, error) {
	return nil, nil
}

func (c *Stream) Close() error {
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

// uint20 =    1048575 =   1gb  = 2b.4bit
// uint22 =    4194303 =   5gb	= 2b.6bit
// uint23 =    8388607 =  11gb  = 2b.7bit
// uint24 =   16777215 =  22gb  = 3b
// uint27 =  134217727 = 177gb  = 3b.3bit
// uint32 = 4294967295 = 5667gb = 4b
// 3b + 3b + 2b.7bit + 2b.7bit = 11b.6bit // 3op + 3op + 3op+1op + 3op+1op
// 4b + 4b + 3b.3bit + 3b.3bit = 14b.6bit // 4op + 4op + 4op + 4op
// 4b + 4b + 4b + 4b           = 16b      // 4op + 4op + 4op + 4op

func signing(b [datagramSigLen]byte, dataEnd int, pkey []byte) [sha256.Size224]byte {
	copy(b[dataEnd+1:sha256.Size224], pkey)
	return sha256.Sum224(b[4 : dataEnd+1+sha256.Size224])
}

func cidFromB(b [datagramSigLen]byte) uint32 {
	return uint32(b[4]) | uint32(b[5])<<8 | uint32(b[6])<<16 | uint32(b[7])<<24
}

func prevDidFromB(b [datagramSigLen]byte) uint32 {
	return uint32(b[8]) | uint32(b[9])<<8 | uint32(b[10])<<16 | uint32(b[11])<<24
}

func didFromB(b [datagramSigLen]byte) uint32 {
	return uint32(b[12]) | uint32(b[13])<<8 | uint32(b[14])<<16 | uint32(b[15])<<24
}
