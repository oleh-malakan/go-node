package node

import (
	"crypto/sha256"
	"net"
	"time"
)

const (
	dataHandshakeBegin = 4
	dataBegin          = 16
	datagramMinLen     = 560
	datagramMaxLen     = 1432
	datagramSigLen     = 1460
)

type incomingDatagram struct {
	b       []byte
	n       int
	dataEnd int
	offset  int
	rAddr   *net.UDPAddr
	next    *incomingDatagram
	err     error
	cid     uint32
	did     uint32
	prevDid uint32
}

func (d *incomingDatagram) checkSig(key []byte) bool {
	copy(d.b[d.n:sha256.Size224], key[:])
	sig := sha256.Sum224(d.b[4 : d.n+sha256.Size224])
	return d.b[0] == sig[0] && d.b[1] == sig[1] && d.b[2] == sig[2] && d.b[3] == sig[3]
}

type outgoingDatagram struct {
	b       [datagramSigLen]byte
	n       int
	dataEnd int
	prev    *outgoingDatagram
	key     []byte
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
		for {
			i := &incomingDatagram{
				b: make([]byte, datagramSigLen),
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
			i.dataEnd = i.n - 1
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

type thread struct{}

type core struct {
	isProcess bool
	next      *core
	nextDrop  chan *core
	drop      chan *core
	signal    chan *struct{}

	inData         chan *incomingDatagram
	lastIncoming   *incomingDatagram
	incoming       *incomingDatagram
	incomingAnchor *incomingDatagram
	outgoing       *outgoingDatagram
	heap           *heap
	/*
		conn              *tls.Conn
		tlsIncomingAnchor *incomingDatagram
		tlsInAnchor       chan *incomingDatagram
		tlsInSignal       chan *struct{}
	*/
}

func (c *core) process() {
	/*
		timerCancelHandshake := time.NewTimer(time.Duration(200) * time.Millisecond)
	*/
	for c.isProcess {
		select {
		case i := <-c.inData:
			if !c.in(i) && c.next != nil {
				c.next.inData <- i

				continue
			}
			/*
				c.tslIn()
			*/
			/*
				case <-c.tlsInSignal:
					c.tslIn()
			*/
		case d := <-c.nextDrop:
			c.next = d.next
			if c.next != nil {
				c.next.drop = c.nextDrop
				select {
				case c.next.signal <- nil:
				default:
				}
			}
			/*
				case <-timerCancelHandshake.C:
					if !c.conn.ConnectionState().HandshakeComplete {
						c.isProcess = false
					}
			*/
		}
	}

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
			if incoming.checkSig(o.key) {
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

/*
func (c *core) tslIn() {
	if c.incomingAnchor != c.lastIncoming {
		select {
		case c.tlsInAnchor <- c.lastIncoming:
			c.incomingAnchor = c.lastIncoming
		default:
		}
	}
}
*/

type tlsConn struct{}

func (c *tlsConn) Read(b []byte) (n int, err error) {
	/*
		var offset int
		for {
			if c.tlsIncomingAnchor != nil {
			LOOP:
				if len(b)-offset <= c.incoming.dataEnd-c.incoming.offset {
					copy(b[offset:], c.incoming.b[c.incoming.offset:c.incoming.offset+len(b)-offset])
					c.incoming.offset = c.incoming.offset + len(b) - offset
					return len(b), nil
				} else {
					copy(b[offset:c.incoming.dataEnd-c.incoming.offset], c.incoming.b[c.incoming.offset:c.lastIncoming.dataEnd])
					offset = offset + c.incoming.dataEnd - c.incoming.offset
					if c.incoming != c.tlsIncomingAnchor {
						c.incoming = c.incoming.next
						goto LOOP
					}
					c.tlsIncomingAnchor = nil
				}
			} else {
				c.tlsInSignal <- nil
				c.tlsIncomingAnchor = <-c.tlsInAnchor
			}
		}
	*/

	return
}

func (c *tlsConn) Write(b []byte) (n int, err error) {
	return
}

func (c *tlsConn) Close() error {
	return nil
}

func (c *tlsConn) LocalAddr() net.Addr {
	return nil
}

func (c *tlsConn) RemoteAddr() net.Addr {
	return nil
}

func (c *tlsConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *tlsConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *tlsConn) SetWriteDeadline(t time.Time) error {
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

type Session struct {
	id [sha256.Size]byte
}

func (s *Session) ID() []byte {
	return s.id[:]
}

func (s *Session) Put(key string, b []byte) error {
	return nil
}

func (s *Session) Get(key string) ([]byte, error) {
	return nil, nil
}

// Sync auto
func (s *Session) Sync() error {
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

func cidFromB(b []byte) uint32 {
	return uint32(b[4]) | uint32(b[5])<<8 | uint32(b[6])<<16 | uint32(b[7])<<24
}

func prevDidFromB(b []byte) uint32 {
	return uint32(b[8]) | uint32(b[9])<<8 | uint32(b[10])<<16 | uint32(b[11])<<24
}

func didFromB(b []byte) uint32 {
	return uint32(b[12]) | uint32(b[13])<<8 | uint32(b[14])<<16 | uint32(b[15])<<24
}

func didFromHandshakeB(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}
