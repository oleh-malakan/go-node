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

	client := &Client{
		controller: &clientController{},
	}

	go client.controller.process()

	return client, nil
}

type Client struct {
	controller *clientController
}

func (c *Client) Connect(nodeID string) (*Connection, error) {
	return &Connection{}, nil
}
