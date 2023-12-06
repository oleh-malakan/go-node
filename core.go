package node

import (
	"crypto/sha256"
	"net"
)

type incomingDatagram struct {
	b      []byte
	n      int
	offset int
	rAddr  *net.UDPAddr
	next   *incomingDatagram
	err    error
}

func (d *incomingDatagram) checkSig(key []byte) bool {
	copy(d.b[d.n:sha256.Size224], key[:])
	sig := sha256.Sum224(d.b[4 : d.n+sha256.Size224])
	return d.b[0] == sig[0] && d.b[1] == sig[1] && d.b[2] == sig[2] && d.b[3] == sig[3]
}

type outgoingDatagram struct {
	b      []byte
	n      int
	offset int
	prev   *outgoingDatagram
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
				b: make([]byte, 1432),
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
}

func (c *core) process() {
	for c.isProcess {
		select {
		case i := <-c.inData:

			//

			if c.next != nil {
				c.next.inData <- i

				continue
			}
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

type Stream struct{}

func (s Stream) NamedStream(id string) (*NamedStream, error) {
	return &NamedStream{}, nil
}

func (s *Stream) Send(b []byte) error {
	return nil
}

func (s *Stream) Receive() ([]byte, error) {
	return nil, nil
}

func (s *Stream) Close() error {
	return nil
}

type NamedStream struct{}

func (s *NamedStream) Send(b []byte) error {
	return nil
}

func (s *NamedStream) Receive() ([]byte, error) {
	return nil, nil
}

func (s *NamedStream) Close() error {
	return nil
}
