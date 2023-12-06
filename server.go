package node

import (
	"crypto/sha256"
	"net"
)

func New(address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) (*Server, error) {
	return &Server{
		nodeAddresses: nodeAddresses,
	}, nil
}

type Server struct {
	address       *net.UDPAddr
	nodeAddresses []*net.UDPAddr
}

func (s *Server) Handler(nodeID string, f func(stream *Stream)) (*Handler, error) {
	return &Handler{
		nodeID: sha256.Sum224([]byte(nodeID)),
	}, nil
}

func (s *Server) Listen(nodeID string) (*Listener, error) {
	return &Listener{
		nodeID: sha256.Sum224([]byte(nodeID)),
	}, nil
}

// clientsLimit default value 524288 if 0
func (s *Server) Run(clientsLimit int) error {
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

func (s *Server) in(c *container, incoming *incomingDatagram) {
	switch {
	case incoming.b[0]&0b10000000 == 0:
		if c.next != nil {
			c.next.inData <- incoming
		}
	case incoming.b[0]&0b10000000 == 1:
		core := &core{
			inData:    make(chan *incomingDatagram),
			nextDrop:  make(chan *core),
			signal:    make(chan *struct{}),
			isProcess: true,
		}
		core.incoming = incoming
		core.lastIncoming = incoming
		core.next = c.next
		c.next = core
		go core.process()
	}
}

func (s *Server) Close() error {
	return nil
}

type Handler struct {
	nodeID [sha256.Size224]byte
}

func (h *Handler) Close() error {
	return nil
}

type Listener struct {
	nodeID [sha256.Size224]byte
}

func (l *Listener) Accept() (*Stream, error) {
	return &Stream{}, nil
}

func (l *Listener) Close() error {
	return nil
}
