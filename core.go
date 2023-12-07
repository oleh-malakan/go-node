package node

import (
	"net"
)

type incomingDatagram struct {
	rAddr  *net.UDPAddr
	next   *incomingDatagram
	b      []byte
	err    error
	n      int
	offset int
}

type outgoingDatagram struct {
	prev   *outgoingDatagram
	b      []byte
	n      int
	offset int
}

type core struct {
	next      *core
	nextDrop  chan *core
	drop      chan *core
	signal    chan *struct{}
	inProcess func(core *core, incoming *incomingDatagram)
	onDestroy func()

	inData       chan *incomingDatagram
	lastIncoming *incomingDatagram
	incoming     *incomingDatagram
	outgoing     *outgoingDatagram
	isProcess    bool
}

func (c *core) process() {
	for c.isProcess {
		select {
		case i := <-c.inData:
			c.inProcess(c, i)
		case d := <-c.nextDrop:
			c.next = d.next
			c.next.drop = c.nextDrop
			c.next.asyncSignal()
		}
	}

	c.onDestroy()
	for {
		select {
		case i := <-c.inData:
			c.next.inData <- i
		case <-c.signal:
		case c.drop <- c:
			return
		}
	}
}

func coreInProcess(core *core, incoming *incomingDatagram) {

	//

	core.next.inData <- incoming
}

func coreEndInProcess(core *core, incoming *incomingDatagram) {}

func coreOnDestroy() {}

func (c *core) asyncSignal() {
	select {
	case c.signal <- nil:
	default:
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

func newCounter() *counter {
	return &counter{
		inc:   make(chan *struct{}),
		dec:   make(chan *struct{}),
		value: make(chan int),
	}
}

type counter struct {
	inc   chan *struct{}
	dec   chan *struct{}
	value chan int
	v     int
}

func (c *counter) process() {
	for {
		select {
		case <-c.inc:
			c.v++
		case <-c.dec:
			c.v--
		case c.value <- c.v:
		}
	}
}
