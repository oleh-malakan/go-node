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
		address:       address,
		nodeAddresses: nodeAddresses,
		controller: &controller{
			tlsConfig: tlsConfig,
			in:        make(chan *incomingPackage, 512),
			nextDrop:  make(chan *container),
		},
	}, nil
}

type Server struct {
	address       *net.UDPAddr
	nodeAddresses []*net.UDPAddr

	controller *controller
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
		incoming := &incomingPackage{
			b: make([]byte, 1432),
		}
		incoming.n, incoming.rAddr, incoming.err = conn.ReadFromUDP(incoming.b)

		s.controller.in <- incoming
	}
}

func (s *Server) Close() error {
	return nil
}
