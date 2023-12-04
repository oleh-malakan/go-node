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

func (s *Server) Handler(nodeID string, f func(stream *Stream)) (*Handler, error) {
	h := &Handler{
		f: f,
	}
	h.nodeID = sha256.Sum224([]byte(nodeID))

	return h, nil
}

func (s *Server) Listen(nodeID string) (*Listener, error) {
	return &Listener{}, nil
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
	case incoming.b[incoming.dataEnd]&0b10000000 == 0:
		incoming.offset = dataBegin
		incoming.cid = cidFromB(incoming.b)
		if c.next != nil {
			c.next.inData <- incoming
		}
	case incoming.b[incoming.dataEnd]&0b10000000 == 1:
		incoming.offset = dataHandshakeBegin
		incoming.did = didFromHandshakeB(incoming.b)
		core := &core{
			heap:        &heap{},
			inData:      make(chan *incomingDatagram),
			nextDrop:    make(chan *core),
			signal:      make(chan *struct{}),
			isProcess:   true,
			tlsInAnchor: make(chan *incomingDatagram),
			tlsInSignal: make(chan *struct{}),
		}
		core.conn = tls.Server(core, s.tlsConfig)
		core.incoming = incoming
		core.lastIncoming = incoming
		core.incomingAnchor = incoming
		core.tlsIncomingAnchor = incoming
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
	f      func(connection *Stream)
}

func (h *Handler) Close() error {
	return nil
}

type Listener struct{}

func (l *Listener) Accept() (*Stream, error) {
	return &Stream{}, nil
}

func (l *Listener) Close() error {
	return nil
}
