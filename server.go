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
