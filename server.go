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
	}, nil
}

type Server struct {
	config        *Config
	tlsConfig     *tls.Config
	address       *net.UDPAddr
	nodeAddresses []*net.UDPAddr
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

	container := &container{
		conn:     conn,
		inData:   make(chan *incomingDatagram),
		nextDrop: make(chan *core),
		in:       s.in,
	}
	container.process()

	return nil
}

func (s *Server) in(c *container, ip *incomingDatagram) {
	switch {
	case ip.b[0]>>7&1 == 0:
		core := &core{
			heap: &heap{
				cap: s.config.HeapCap,
			},
			inData:   make(chan *incomingDatagram),
			nextDrop: make(chan *core),
			reset:    make(chan *struct{}),
		}
		core.conn = tls.Server(core, s.tlsConfig)
		core.incoming = ip
		core.lastIncoming = ip
		core.next = c.next
		c.next = core
		go core.process()
	case ip.b[0]>>7&1 == 1:
		ip.cid = bToID(ip.b[4:7])
		if c.next != nil {
			c.next.inData <- ip
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
