package node

import (
	"crypto/tls"
	"errors"
	"net"
)

type Connection struct{}

func (c *Connection) Send(b []byte) error {
	return nil
}

func (c *Connection) Receive() ([]byte, error) {
	return nil, nil
}

func (c *Connection) Close() error {
	return nil
}

func Handler(nodeID string, f func(query []byte, connection *Connection)) {}

func Do(address *net.UDPAddr, nodeAddress ...*net.UDPAddr) error {
	return nil
}

func DoTLS(tlsConfig *tls.Config, address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) error {
	return nil
}

type Client struct{}

func (c *Client) Connect(nodeID string, query []byte) (*Connection, error) {
	return &Connection{}, nil
}

func Dial(nodeAddresses ...*net.UDPAddr) (*Client, error) {
	if len(nodeAddresses) == 0 {
		return nil, errors.New("node address not specified")
	}

	return &Client{}, nil
}

func DialTLS(tlsConfig *tls.Config, nodeAddresses ...*net.UDPAddr) (*Client, error) {
	if len(nodeAddresses) == 0 {
		return nil, errors.New("node address not specified")
	}

	return &Client{}, nil
}
