package node

import (
	"crypto/ecdh"
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
