package node

import (
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"net"
)

func New(tlsConfig *tls.Config, address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) (*Server, error) {
	if tlsConfig == nil {
		return nil, errors.New("require tls config")
	}

	server := &Server{
		tlsConfig:  tlsConfig,
		memoryLock: make(chan *struct{}, 1),
		checkLock:  make(chan *struct{}, 1),
	}
	if err := server.do(tlsConfig, address, nodeAddresses...); err != nil {
		return nil, err
	}

	return server, nil
}

type client struct {
	conn         *tls.Conn
	lastReadData *readData
	readData     *readData
	writeData    *writeData
	initalMac    [32]byte
	heap         *heap
	next         *client
	lock         chan *struct{}
	drop         bool
}

type Server struct {
	handlers   []*Handler
	memory     *client
	memoryLock chan *struct{}
	checkLock  chan *struct{}
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
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		return err
	}

	for {
		r := &readData{
			b: make([]byte, 1432),
		}
		r.n, r.rAddr, r.err = conn.ReadFromUDP(r.b)

		go n.process(r)
	}
}

func (n *Server) process(r *readData) {
	switch {
	case r.b[0]>>7&1 == 0:
		r.nextMac = sha256.Sum256(r.b[1:r.n])
		n.checkLock <- nil
		var current *client
		if !n.bypass(func(c *client) bool {
			return compareID(c.initalMac[0:32], r.nextMac[0:32])
		}) {
			current = &client{
				conn: tls.Server(&dataport{}, n.tlsConfig),
				lock: make(chan *struct{}, 1),
				heap: &heap{
					cap: 512,
				},
			}
			current.readData = r
			current.lastReadData = r
			current.initalMac = r.nextMac
			current.drop = false
		}
		n.memoryLock <- nil
		if current != nil {
			if n.memory != nil {
				current.next = n.memory
				n.memory = current
			} else {
				n.memory = current
			}
		}
		<-n.memoryLock
		<-n.checkLock
	case r.b[0]>>7&1 == 1 && n.memory != nil:
		r.nextMac = sha256.Sum256(r.b[65:r.n])
		n.bypass(func(c *client) (f bool) {
			w := c.writeData
			for w != nil {
				if compareID(w.mac[0:32], r.b[1:33]) {
					if compareID(c.lastReadData.nextMac[0:32], r.b[33:65]) {
						c.lastReadData.next = r
						c.lastReadData = r
						var last *readData
						c.lastReadData.next, last = c.heap.find(r.nextMac[0:32])
						if last != nil {
							c.lastReadData = last
						}
					} else {
						c.heap.put(r)
					}

					return true
				}
				w = w.prev
			}

			return
		})
	}
}

func (n *Server) bypass(f func(c *client) bool) (ok bool) {
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
			if ok = f(current); ok {
				next = nil
				goto UNLOCK
			}
		}
		next = current.next
		if next != nil {
			next.lock <- nil
		}
	UNLOCK:
		<-current.lock
		current = next
	}

	return
}

func (n *Server) Run() {

}
