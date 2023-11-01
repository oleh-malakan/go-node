package node

import (
	"errors"
	"net"
)

type Client struct {
	cConn        chan *net.UDPConn
	cReadConnErr chan error

	cErr chan error
}

func (c *Client) Connect(nodeID string, query []byte) (*Connection, error) {
	return &Connection{}, nil
}

func Dial(nodeAddresses ...*net.UDPAddr) (*Client, error) {
	if len(nodeAddresses) == 0 {
		return nil, errors.New("node address not specified")
	}

	client := &Client{}
	client.dial(nodeAddresses...)

	return client, nil
}

func (c *Client) dial(nodeAddresses ...*net.UDPAddr) {
	c.cConn = make(chan *net.UDPConn)
	c.cReadConnErr = make(chan error)

	go func() {
		var conn *net.UDPConn
	NEW:
		var err error
		conn, err = net.DialUDP("udp", nil, nodeAddresses[0])
		if err != nil {
			c.cErr <- err
			goto NEW
		}

		for {
			select {
			case c.cConn <- conn:
			case <-c.cReadConnErr:
				goto NEW
			}
		}
	}()

	go c.runRead()
	go c.runWrite()
}

func (c *Client) runRead() {}

func (c *Client) runWrite() {}
