package node

import (
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
	nodeAddresses []*net.UDPAddr
}

func (c *Client) Connect(nodeID string) (*Stream, error) {
	return &Stream{}, nil
}

func (c *Client) process() {
	conn, err := net.DialUDP("udp", nil, c.nodeAddresses[0])
	if err != nil {

	}

	currentCore := &core{
		inData:         make(chan *incomingDatagram),
		drop:           make(chan *core, 1),
		isProcess:      true,
		inProcess:      coreInProcess,
		destroyProcess: coreDestroyProcess,
	}
	endCore := &core{
		inData:         make(chan *incomingDatagram),
		drop:           make(chan *core),
		isProcess:      true,
		inProcess:      coreEndInProcess,
		destroyProcess: coreDestroyProcess,
	}

	currentCore.next = endCore

	go currentCore.process()
	go endCore.process()

	for {
		i := &incomingDatagram{
			b: make([]byte, 1432),
		}
		i.n, i.rAddr, i.err = conn.ReadFromUDP(i.b)
		if i.err != nil {

			//continue
		}
		currentCore.inData <- i
	}
}
