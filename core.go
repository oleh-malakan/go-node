package node

import (
	"crypto/cipher"
	"net"
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
	did    int64
}

type cIDDatagram struct {
	datagram *datagram
	cid      int32
}

func parseCIDDatagram(d *datagram) *cIDDatagram {
	return &cIDDatagram{
		datagram: d,
	}
}

type core struct {
	inData       chan *datagram
	drop         chan *core
	lastIncoming *datagram
	incoming     *datagram
	cid          int32
	aead         cipher.AEAD
	isProcess    bool
}

func (c *core) process() {
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
