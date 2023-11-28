package node

import (
	"crypto/tls"
	"errors"
	"net"
)

func Dial(tlsConfig *tls.Config, nodeAddresses ...*net.UDPAddr) (*Client, error) {
	if len(nodeAddresses) == 0 {
		return nil, errors.New("node address not specified")
	}

	client := &Client{}

	go client.process()

	return client, nil
}

type Client struct {
}

func (c *Client) Connect(nodeID string) (*Connection, error) {
	return &Connection{}, nil
}

func (c *Client) process() {

}

type clientContainer struct{}

func (c *clientContainer) process() {

}
