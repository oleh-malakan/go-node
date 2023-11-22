package node

import (
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"net"
)

func Run(tlsConfig *tls.Config, address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) (*Server, error) {
	if tlsConfig == nil {
		return nil, errors.New("require tls config")
	}

	node := &Server{}
	if err := node.do(tlsConfig, address, nodeAddresses...); err != nil {
		return nil, err
	}

	return node, nil
}

type client struct {
	conn         *tls.Conn
	lastReadData *readData
	readData     *readData
	heap         *heap
	writeData    *writeData
	next         *client
	lock         chan *struct{}
	drop         bool
}

type Server struct {
	handlers   []*Handler
	memory     *client
	memoryLock chan *struct{}
	tlsConfig  *tls.Config
}

func (n *Server) Handler(nodeID string, f func(connection *Connection)) (*Handler, error) {
	h := &Handler{
		f: f,
	}
	h.nodeID = sha256.Sum256([]byte(nodeID))

	n.handlers = append(n.handlers, h)

	return h, nil
}

func (n *Server) Listen(nodeID string) (*Listener, error) {
	return &Listener{}, nil
}

func (n *Server) do(tlsConfig *tls.Config, address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) error {
	n.tlsConfig = tlsConfig

	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		return err
	}

	n.memoryLock = make(chan *struct{}, 1)
	for {
		r := &readData{
			b: make([]byte, 1432),
		}
		r.n, r.rAddr, r.err = conn.ReadFromUDP(r.b)

		go n.bypass(r)
	}
}

func (n *Server) bypass(r *readData) {
	switch {
	case r.b[0]>>7&1 == 0:
		r.nextMac = sha256.Sum256(r.b[1:r.n])
		client := &client{
			conn: tls.Server(&dataport{}, n.tlsConfig),
			lock: make(chan *struct{}, 1),
			heap: &heap{
				cap: 512,
			},
		}
		client.readData = r
		client.lastReadData = r
		client.drop = false

		n.memoryLock <- nil
		if n.memory != nil {
			client.next = n.memory
			n.memory = client
		} else {
			n.memory = client
		}
		<-n.memoryLock
	case r.b[0]>>7&1 == 1 && n.memory != nil:
		r.nextMac = sha256.Sum256(r.b[65:r.n])
		var current *client
		n.memoryLock <- nil
		if n.memory != nil {
			current = n.memory
			current.lock <- nil
		}
		<-n.memoryLock
		for current != nil {
			var next *client
			if !current.drop {
				w := current.writeData
				for w != nil {
					if compareID(w.mac[0:32], r.b[1:33]) {
						if compareID(current.lastReadData.nextMac[0:32], r.b[33:65]) {
							current.lastReadData.next = r
							current.lastReadData = r
							var last *readData
							current.lastReadData.next, last = current.heap.find(r.nextMac[0:32])
							if last != nil {
								current.lastReadData = last
							}
						} else {
							current.heap.put(r)
						}

						next = nil
						goto FOUND
					}
					w = w.prev
				}
			}

			next = current.next
			if next != nil {
				next.lock <- nil
			}
		FOUND:
			<-current.lock
			current = next
		}
	}
}

func (n *Server) Wait() {

}
