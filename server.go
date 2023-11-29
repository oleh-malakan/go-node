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
		config:        &config,
		tlsConfig:     tlsConfig,
		address:       address,
		nodeAddresses: nodeAddresses,
		in:            make(chan *incomingPackage),
		nextDrop:      make(chan *serverContainer),
	}, nil
}

type Server struct {
	config        *Config
	tlsConfig     *tls.Config
	address       *net.UDPAddr
	nodeAddresses []*net.UDPAddr

	next     *serverContainer
	in       chan *incomingPackage
	nextDrop chan *serverContainer
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

	go s.process()

	for {
		i := &incomingPackage{
			b: make([]byte, 1460),
		}
		i.n, i.rAddr, i.err = conn.ReadFromUDP(i.b)

		s.in <- i
	}
}

func (s *Server) process() {
	for {
		select {
		case i := <-s.in:
			switch {
			case i.b[0]>>7&1 == 0:
				new := &serverContainer{
					core: &core{
						heap: &heap{
							cap: s.config.HeapCap,
						},
					},
					in:       make(chan *incomingPackage),
					nextDrop: make(chan *serverContainer),
					reset:    make(chan *struct{}),
				}
				new.core.conn = tls.Server(new.core, s.tlsConfig)
				new.core.incoming = i
				new.core.lastIncoming = i
				new.next = s.next
				s.next = new
				go new.process()
			case i.b[0]>>7&1 == 1:
				i.cid = cid(i.b[4:10])
				if s.next != nil {
					s.next.in <- i
				}
			}
		case d := <-s.nextDrop:
			s.next = d.next
			if s.next != nil {
				s.next.drop = s.nextDrop
				select {
				case s.next.reset <- nil:
				default:
				}
			}
		}
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
			if !c.core.in(i) && c.next != nil {
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
