package node

import (
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"net"
)

func New(tlsConfig *tls.Config, nodeAddresses ...*net.UDPAddr) (*Server, error) {
	if tlsConfig == nil {
		return nil, errors.New("require tls config")
	}

	return &Server{
		tlsConfig:     tlsConfig,
		nodeAddresses: nodeAddresses,
	}, nil
}

type Server struct {
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

func (s *Server) HandshakeProcess() {

}

// clientsLimit default value 524288 if 0
func (s *Server) Run(address *net.UDPAddr, clientsLimit int) error {
	if clientsLimit <= 0 {
		clientsLimit = 524288
	}

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
			heap:     &heap{},
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
