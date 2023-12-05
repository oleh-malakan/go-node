package node

import (
	"errors"
	"net"
)

func Connect(nodeAddresses ...*net.UDPAddr) (*Client, error) {
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
	nodeAddresses []*net.UDPAddr
}

func (c *Client) Stream(nodeID string) (*Stream, error) {
	return &Stream{}, nil
}

func (c *Client) process() {
	conn, err := net.DialUDP("udp", nil, c.nodeAddresses[0])
	if err != nil {

	}

	container := &container{
		conn:     conn,
		inData:   make(chan *incomingDatagram),
		nextDrop: make(chan *core),
		in:       c.in,
		next: &core{
			heap:      &heap{},
			inData:    make(chan *incomingDatagram),
			nextDrop:  make(chan *core),
			signal:    make(chan *struct{}),
			isProcess: true,
		},
	}
	go container.next.process()
	container.process()
}

func (client *Client) in(c *container, incoming *incomingDatagram) {
	incoming.offset = dataBegin
	incoming.cid = cidFromB(incoming.b)
	c.next.inData <- incoming
}
