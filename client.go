package node

import (
	"crypto/tls"
	"net"
)

type Client struct{}

func (c *Client) Connect(nodeID string, query []byte) (*Connection, error) {
	return &Connection{}, nil
}

func Dial(tlsConfig *tls.Config, nodeAddresses ...*net.UDPAddr) (*Client, error) {
	if len(nodeAddresses) == 0 {
		return nil, newError("node address not specified")
	}

	client := &Client{}
	err := client.dial(tlsConfig, nodeAddresses...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) dial(tlsConfig *tls.Config, nodeAddresses ...*net.UDPAddr) error {
	conn, err := net.DialUDP("udp", nil, nodeAddresses[0])
	if err != nil {
		return err
	}

	return nil
}
