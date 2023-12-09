package node

import (
	"crypto/cipher"
	"encoding/binary"
	"net"
)

const (
	datagramMinLen = 560
	datagramCap    = 1432
	dataCap        = 1416
)

type incomingDatagram struct {
	rAddr   *net.UDPAddr
	next    *incomingDatagram
	cipherB []byte
	b       []byte
	cid     *id
	err     error
	n       int
	offset  int
}

func (d *incomingDatagram) prepareCID() {
	d.cid = &id{
		ID1: binary.BigEndian.Uint64(d.cipherB[1:9]),
		ID2: binary.BigEndian.Uint64(d.cipherB[9:17]),
		ID3: binary.BigEndian.Uint64(d.cipherB[17:25]),
		ID4: binary.BigEndian.Uint64(d.cipherB[25:33]),
	}
}

func (d *incomingDatagram) decode() bool {
	return true
}

type outgoingDatagram struct {
	prev   *outgoingDatagram
	b      []byte
	n      int
	offset int
}

type id struct {
	ID1 uint64
	ID2 uint64
	ID3 uint64
	ID4 uint64
}

type core struct {
	next           *core
	drop           chan *core
	inProcess      func(core *core, incoming *incomingDatagram)
	destroyProcess func()

	inData       chan *incomingDatagram
	lastIncoming *incomingDatagram
	incoming     *incomingDatagram
	outgoing     *outgoingDatagram
	cid          *id
	aead         cipher.AEAD
	isProcess    bool
}

func (c *core) process() {
	for c.isProcess {
		select {
		case i := <-c.inData:
			c.inProcess(c, i)
		case d := <-c.next.drop:
			c.next = d.next
		}
	}

	c.destroyProcess()
	for {
		select {
		case i := <-c.inData:
			c.next.inData <- i
		case c.drop <- c:
			return
		}
	}
}

func coreInProcess(core *core, incoming *incomingDatagram) {
	if incoming.decode() {

		//

	}
}

func coreEndInProcess(core *core, incoming *incomingDatagram) {}

func coreDestroyProcess() {}

type Stream struct{}

func (s Stream) MakeStream(id string) (*NamedStream, error) {
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

type transport struct {
	conn *net.UDPConn
}

func (t *transport) write(b []byte, addr *net.UDPAddr) (int, error) {
	return t.conn.WriteToUDP(b, addr)
}

func (t *transport) read(b []byte) (int, *net.UDPAddr, error) {
	return t.conn.ReadFromUDP(b)
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
