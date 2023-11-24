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
		tlsConfig:     tlsConfig,
		address:       address,
		nodeAddresses: nodeAddresses,
		memory: &memory{
			tlsConfig: tlsConfig,
			in:        make(chan *incomingPackage, 512),
			nextDrop:  make(chan *node),
		},
	}
	go server.memory.process()

	return server, nil
}

type Server struct {
	tlsConfig     *tls.Config
	memory        *memory
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

	for {
		incoming := &incomingPackage{
			b: make([]byte, 1432),
		}
		incoming.n, incoming.rAddr, incoming.err = conn.ReadFromUDP(incoming.b)

		s.memory.in <- incoming
	}
}

func (s *Server) Close() error {
	return nil
}
