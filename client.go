package node

import (
	"crypto/tls"
	"errors"
	"net"
)

func Dial(tlsConfig *tls.Config, nodeAddresses ...*net.UDPAddr) (*Client, error) {
	if tlsConfig == nil {
		return nil, errors.New("require tls config")
	}

	if len(nodeAddresses) == 0 {
		return nil, errors.New("node address not specified")
	}

	client := &Client{
		tlsConfig: tlsConfig,
	}

	go client.process()

	return client, nil
}

type Client struct {
	tlsConfig *tls.Config
	newConn   chan *net.UDPAddr
	next      *core
	nextDrop  chan *core
}

func (c *Client) Connect(nodeID string) (*Connection, error) {
	return &Connection{}, nil
}

func (c *Client) process() {
	for {
		select {
		case addr := <-c.newConn:
			conn, err := net.DialUDP("udp", nil, addr)
			if err != nil {
				continue
			}
			core := &core{
				heap: &heap{
					cap: 512,
				},
				in:       make(chan *incomingPackage),
				nextDrop: make(chan *core),
				reset:    make(chan *struct{}),
			}
			core.conn = tls.Server(core, c.tlsConfig)
			core.next = c.next
			c.next = core
			go core.process()

			reader := &reader{
				conn: conn,
				core: core,
			}
			core.reader = reader
			go reader.process()

		case d := <-c.nextDrop:
			c.next = d.next
			if c.next != nil {
				c.next.drop = c.nextDrop
				select {
				case c.next.reset <- nil:
				default:
				}
			}
		}
	}
}
