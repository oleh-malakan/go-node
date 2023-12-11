package node

import (
	"crypto/cipher"
	"crypto/ecdh"
	"net"

	"github.com/oleh-malakan/go-node/internal"
)

const (
	datagramMinLen = 560
	datagramCap    = 1432
	dataCap        = 1416
)

type datagram struct {
	rAddr  *net.UDPAddr
	next   *datagram
	b      []byte
	data   []byte
	err    error
	n      int
	begin  int
	offset int
}

type cIDDatagram struct {
	datagram *datagram
	cid      int64
}

func parseCIDDatagram(d *datagram) *cIDDatagram {
	return &cIDDatagram{
		datagram: d,
	}
}

type core struct {
	inData         chan *datagram
	drop           chan int
	lastIncoming   *datagram
	incoming       *datagram
	privateKey     *ecdh.PrivateKey
	publicKeyBytes []byte
	cid            int64
	aead           cipher.AEAD
	isProcess      bool
}

func (c *core) process() {
	// send core.index

	for c.isProcess {
		select {
		case <-c.inData:
		}
	}
}

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

type controller struct {
	memory           *internal.IndexArray[core]
	drop             chan int
	counter          *counter
	connectionsLimit int
}

func (c *controller) in(incoming *datagram) {
	cIDDatagram := parseCIDDatagram(incoming)
	if current := c.memory.Get(cIDDatagram.cid); current != nil && current.check() {
		current.inData <- incoming
	} else {
		if <-c.counter.value < c.connectionsLimit {
			c.counter.inc <- nil

			go func() {

				new := &core{
					inData:       make(chan *datagram),
					drop:         c.drop,
					isProcess:    true,
					incoming:     incoming,
					lastIncoming: incoming,
				}
				new.cid = c.memory.Put(new)
				go new.process()
			}()
		}
	}
}


func (c *controller) free(index int64) {
	c.memory.Free(index)
	c.counter.dec <- nil
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
