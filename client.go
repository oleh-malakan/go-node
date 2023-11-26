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
	client.conn = tls.Client(&dataport{}, tlsConfig)
	err := client.dial(nodeAddresses...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

type Client struct {
	conn *tls.Conn
}

func (c *Client) Connect(nodeID string) (*Connection, error) {
	return &Connection{}, nil
}

func (c *Client) dial(nodeAddresses ...*net.UDPAddr) error {
	/*
		conn, err := net.DialUDP("udp", nil, nodeAddresses[0])
		if err != nil {
			return err
		}

		go c.handshake()

		go func() {
			for {
				select {
				case b := <-c.cWrite:

				}
			}

		}()
	*/
	return nil
}

func (c *Client) handshake() {
	if err := c.conn.Handshake(); err != nil {

	}
}
