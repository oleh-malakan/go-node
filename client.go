package node

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"errors"
	"net"
)

func Dial(nodeAddresses ...*net.UDPAddr) (*Client, error) {
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
	nodeAddresses   []*net.UDPAddr
	serverPublicKey *ecdh.PublicKey
	transport       *transport
}

func (c *Client) Connect(nodeID string) (*Stream, error) {
	return &Stream{}, nil
}

func (c *Client) Close() error {
	return nil
}

func (c *Client) process() {
	conn, err := net.DialUDP("udp", nil, c.nodeAddresses[0])
	if err != nil {

	}

	c.transport = &transport{
		conn: conn,
	}

	inData := make(chan *datagram)

	go func() {
		privateKey, err := ecdh.X25519().GenerateKey(rand.Reader)
		if err != nil {

		}

		b := make([]byte, datagramMinLen)
		copy(b[1:33], privateKey.PublicKey().Bytes())

		_, err = c.transport.writeUDP(b, c.nodeAddresses[0])
		if err != nil {

		}

		for {
			select {
			case i := <-inData:
				remotePublicKey, err := ecdh.X25519().NewPublicKey(i.b[1:33])
				if err != nil {
					continue
				}
				secret, err := privateKey.ECDH(remotePublicKey)
				if err != nil {
					continue
				}
				block, err := aes.NewCipher(secret)
				if err != nil {
					continue
				}
				cipher, err := cipher.NewGCM(block)
				if err != nil {
					continue
				}
				cipher.Open()

				core := &core{
					inData:    inData,
					isProcess: true,
					cipher:    cipher,
				}
				go core.process()

				return
			}
		}
	}()

	for {
		i := &datagram{
			b:     make([]byte, 1432),
			rAddr: c.nodeAddresses[0],
		}
		i.n, i.err = c.transport.read(i.b)
		if i.err != nil {

			//continue
		}
		inData <- i
	}
}
