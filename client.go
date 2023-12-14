package node

import (
	"crypto/ecdh"
	"errors"
	"net"
)

func NewClient(nodeAddresses ...*net.UDPAddr) (*Client, error) {
	client, err := newClient(nodeAddresses...)
	if err != nil {
		return nil, err
	}

	go client.process()

	return client, nil
}

func NewClientPrivateHello(serverPublicKey []byte, nodeAddresses ...*net.UDPAddr) (*Client, error) {
	client, err := newClient(nodeAddresses...)
	if err != nil {
		return nil, err
	}

	client.serverPublicKey, err = ecdh.X25519().NewPublicKey(serverPublicKey)
	if err != nil {
		return nil, err
	}

	go client.process()

	return client, nil
}

func newClient(nodeAddresses ...*net.UDPAddr) (*Client, error) {
	if len(nodeAddresses) == 0 {
		return nil, errors.New("node address not specified")
	}

	client := &Client{
		nodeAddresses: nodeAddresses,
	}

	return client, nil
}

type Client struct {
	nodeAddresses   []*net.UDPAddr
	serverPublicKey *ecdh.PublicKey
}

func (c *Client) Connect(nodeID string) (*Stream, error) {
	return &Stream{}, nil
}

func (c *Client) process() {
	conn, err := net.DialUDP("udp", nil, c.nodeAddresses[0])
	if err != nil {

	}

	core := &core{
		inData:    make(chan *datagram),
		isProcess: true,
	}

	go core.process()

	for {
		i := &datagram{
			b: make([]byte, 1432),
		}
		i.n, i.rAddr, i.err = conn.ReadFromUDP(i.b)
		if i.err != nil {

			//continue
		}
		core.inData <- i
	}
}
