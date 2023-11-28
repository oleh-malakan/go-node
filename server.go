package node

import (
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"net"
)

type Config struct {
	ClientsLimit int // default value 524288 if 0
	HeapCap      int // default value 512    if 0
}

func New(config Config, tlsConfig *tls.Config, address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) (*Server, error) {
	if tlsConfig == nil {
		return nil, errors.New("require tls config")
	}

	if config.ClientsLimit <= 0 {
		config.ClientsLimit = 524288
	}
	if config.HeapCap <= 0 {
		config.HeapCap = 512
	}

	return &Server{
		address:       address,
		nodeAddresses: nodeAddresses,
		controller: &serverController{
			config:    &config,
			tlsConfig: tlsConfig,
			in:        make(chan *incomingPackage),
			nextDrop:  make(chan *serverContainer),
		},
	}, nil
}

type Server struct {
	address       *net.UDPAddr
	nodeAddresses []*net.UDPAddr

	controller *serverController
}

func (s *Server) Handler(nodeID string, f func(connection *Connection)) (*Handler, error) {
	h := &Handler{
		f: f,
	}
	h.nodeID = sha256.Sum256([]byte(nodeID))

	return h, nil
}

func (s *Server) Listen(nodeID string) (*Listener, error) {
	return &Listener{}, nil
}

func (s *Server) Run() error {
	conn, err := net.ListenUDP("udp", s.address)
	if err != nil {
		return err
	}

	go s.controller.process()

	for {
		i := &incomingPackage{
			b: make([]byte, 1432),
		}
		i.n, i.rAddr, i.err = conn.ReadFromUDP(i.b)

		s.controller.in <- i
	}
}

func (s *Server) Close() error {
	return nil
}

type Handler struct {
	nodeID [32]byte
	f      func(connection *Connection)
}

func (h *Handler) Close() error {
	return nil
}

type Listener struct{}

func (l *Listener) Accept() (*Connection, error) {
	return &Connection{}, nil
}

func (l *Listener) Close() error {
	return nil
}

type serverController struct {
	config    *Config
	tlsConfig *tls.Config
	next      *serverContainer
	in        chan *incomingPackage
	nextDrop  chan *serverContainer
}

func (c *serverController) process() {
	for {
		select {
		case i := <-c.in:
			switch {
			case i.b[0]>>7&1 == 0:
				new := &serverContainer{
					core: &core{
						iPKey: sha256.Sum256(i.b[1:i.n]),
						heap: &heap{
							cap: c.config.HeapCap,
						},
					},
					in:       make(chan *incomingPackage),
					nextDrop: make(chan *serverContainer),
					reset:    make(chan *struct{}),
				}
				new.core.conn = tls.Server(new.core, c.tlsConfig)
				new.core.incoming = i
				new.core.lastIncoming = i
				new.next = c.next
				c.next = new
				go new.process()
			case i.b[0]>>7&1 == 1:
				if c.next != nil {
					c.next.in <- i
				}
			}
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

type serverContainer struct {
	core     *core
	next     *serverContainer
	in       chan *incomingPackage
	nextDrop chan *serverContainer
	drop     chan *serverContainer
	reset    chan *struct{}
	isDrop   bool
}

func (c *serverContainer) process() {
	for !c.isDrop {
		select {
		case i := <-c.in:
			o := c.core.outgoing
			for o != nil {
				if compare16(o.b[1:17], i.b[1:17]) {
					c.core.in(i)
					continue
				}
				o = o.prev
			}

			if c.next != nil {
				c.next.in <- i
			}
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

	for {
		select {
		case i := <-c.in:
			if c.next != nil {
				c.next.in <- i
			}
		case <-c.reset:
		case c.drop <- c:
			return
		}
	}
}
