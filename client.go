package node

import (
	"crypto/ecdh"
	"crypto/rand"
	"errors"
	"net"
)

func Dial(nodeAddresses ...*net.UDPAddr) (*Client, error) {
	if len(nodeAddresses) == 0 {
		return nil, errors.New("node address not specified")
	}

	client := &Client{
		nodeAddresses: nodeAddresses,
	}

	go client.process()

	return client, nil
}

type Client struct {
	nodeAddresses   []*net.UDPAddr
	serverPublicKey *ecdh.PublicKey
	transport       *transport
}

func (c *Client) Connect(nodeID string) (*Stream, error) {
	return &Stream{}, nil
}

func (c *Client) Close() error {
	return nil
}

func (c *Client) process() {
	conn, err := net.DialUDP("udp", nil, c.nodeAddresses[0])
	if err != nil {

	}

	c.transport = &transport{
		conn: conn,
	}

	privateKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {

	}

	b := make([]byte, datagramMinLen)
	copy(b[1:33], privateKey.PublicKey().Bytes())

	_, err = c.transport.writeUDP(b, c.nodeAddresses[0])
	if err != nil {

	}

	_, err = c.transport.read(b)
	if err != nil {

	}

	core := &core{
		inData:    make(chan *datagram),
		isProcess: true,
	}

	go core.process()

	for {
		i := &datagram{
			b:     make([]byte, 1432),
			rAddr: c.nodeAddresses[0],
		}
		i.n, i.err = c.transport.read(i.b)
		if i.err != nil {

			//continue
		}
		core.inData <- i
	}
}
