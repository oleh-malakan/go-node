package node

import (
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"net"
	"sync"
)

func New(tlsConfig *tls.Config, address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) (*Server, error) {
	if tlsConfig == nil {
		return nil, errors.New("require tls config")
	}

	return &Server{
		tlsConfig:     tlsConfig,
		address:       address,
		nodeAddresses: nodeAddresses,
	}, nil
}

type Server struct {
	handlers      []*Handler
	memory        *client
	memoryLock    sync.Mutex
	checkLock     sync.Mutex
	tlsConfig     *tls.Config
	address       *net.UDPAddr
	nodeAddresses []*net.UDPAddr
}

type client struct {
	conn         *tls.Conn
	lastReadData *incomingPackage
	readData     *incomingPackage
	writeData    *outgoingPackage
	initalMac    [32]byte
	heap         *heap
	next         *client
	lock         sync.Mutex
	drop         bool
}

func (s *Server) Handler(nodeID string, f func(connection *Connection)) (*Handler, error) {
	h := &Handler{
		f: f,
	}
	h.nodeID = sha256.Sum256([]byte(nodeID))

	s.handlers = append(s.handlers, h)

	return h, nil
}

func (s *Server) Listen(nodeID string) (*Listener, error) {
	return &Listener{}, nil
}

func (s *Server) process(p *incomingPackage) {
	switch {
	case p.b[0]>>7&1 == 0:
		p.nextMac = sha256.Sum256(p.b[1:p.n])
		s.checkLock.Lock()
		var current *client
		if !s.bypass(func(c *client) bool {
			return compareID(c.initalMac[0:32], p.nextMac[0:32])
		}) {
			current = &client{
				conn: tls.Server(&dataport{}, s.tlsConfig),
				heap: &heap{
					cap: 512,
				},
			}
			current.readData = p
			current.lastReadData = p
			current.initalMac = p.nextMac
			current.drop = false
		}
		s.memoryLock.Lock()
		if current != nil {
			if s.memory != nil {
				current.next = s.memory
				s.memory = current
			} else {
				s.memory = current
			}
		}
		s.memoryLock.Unlock()
		s.checkLock.Unlock()
	case p.b[0]>>7&1 == 1 && s.memory != nil:
		p.nextMac = sha256.Sum256(p.b[65:p.n])
		s.bypass(func(c *client) (f bool) {
			w := c.writeData
			for w != nil {
				if compareID(w.mac[0:32], p.b[1:33]) {
					if compareID(c.lastReadData.nextMac[0:32], p.b[33:65]) {
						c.lastReadData.next = p
						c.lastReadData = p
						var last *incomingPackage
						c.lastReadData.next, last = c.heap.find(p.nextMac[0:32])
						if last != nil {
							c.lastReadData = last
						}
					} else {
						c.heap.put(p)
					}

					return true
				}
				w = w.prev
			}

			return
		})
	}
}

func (s *Server) bypass(f func(c *client) bool) (ok bool) {
	var current *client
	s.memoryLock.Lock()
	if s.memory != nil {
		current = s.memory
		current.lock.Lock()
	}
	s.memoryLock.Unlock()
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
			next.lock.Lock()
		}
	UNLOCK:
		current.lock.Unlock()
		current = next
	}

	return
}

func (s *Server) Run() error {
	conn, err := net.ListenUDP("udp", s.address)
	if err != nil {
		return err
	}

	for {
		r := &incomingPackage{
			b: make([]byte, 1432),
		}
		r.n, r.rAddr, r.err = conn.ReadFromUDP(r.b)

		go s.process(r)
	}
}

func (s *Server) Close() error {
	return nil
}
