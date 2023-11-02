package node

import (
	"errors"
	"net"
)

type Client struct {
	conn *net.UDPConn
}

func (c *Client) Connect(nodeID string, query []byte) (*Connection, error) {
	return &Connection{}, nil
}

func Dial(nodeAddresses ...*net.UDPAddr) (*Client, error) {
	if len(nodeAddresses) == 0 {
		return nil, errors.New("node address not specified")
	}

	client := &Client{}
	err := client.dial(nodeAddresses...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) dial(nodeAddresses ...*net.UDPAddr) (err error) {	
	c.conn, err = net.DialUDP("udp", nil, nodeAddresses[0])
	if err != nil {
		return
	}

	c.runRead()
	c.runWrite()

	return
}

func (c *Client) runRead() {
	go func() {

	}()
}

func (c *Client) runWrite() {
	go func() {

	}()
}
