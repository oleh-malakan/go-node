package node

import (
	"crypto/tls"
	"errors"
	"net"
)

func Connect(tlsConfig *tls.Config, nodeAddresses ...*net.UDPAddr) (*Client, error) {
	if tlsConfig == nil {
		return nil, errors.New("require tls config")
	}

	if len(nodeAddresses) == 0 {
		return nil, errors.New("node address not specified")
	}

	client := &Client{
		tlsConfig:     tlsConfig,
		nodeAddresses: nodeAddresses,
	}

	go client.process()

	return client, nil
}

type Client struct {
	tlsConfig     *tls.Config
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
			heap:     &heap{},
			inData:   make(chan *incomingDatagram),
			nextDrop: make(chan *core),
			reset:    make(chan *struct{}),
		},
	}
	container.next.conn = tls.Server(container.next, c.tlsConfig)
	go container.next.process()
	container.process()
}

func (client *Client) in(c *container, ip *incomingDatagram) {
	ip.cid = bToID(ip.b[cidBegin:cidEnd])
	c.next.inData <- ip
}
