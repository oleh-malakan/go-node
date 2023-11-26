package node

import (
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"net"
)

type Config struct {
	BufferSize   int // default value 65536  if 0
	ClientsLimit int // default value 524288 if 0
	HeapSize     int // default value 512    if 0
}

func New(config Config, tlsConfig *tls.Config, address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) (*Server, error) {
	if tlsConfig == nil {
		return nil, errors.New("require tls config")
	}

	if config.BufferSize <= 0 {
		config.BufferSize = 65536
	}
	if config.ClientsLimit <= 0 {
		config.ClientsLimit = 524288
	}
	if config.HeapSize <= 0 {
		config.HeapSize = 512
	}

	return &Server{
		config:        &config,
		address:       address,
		nodeAddresses: nodeAddresses,
		controller: &controller{
			tlsConfig: tlsConfig,
			in:        make(chan *incomingPackage, config.BufferSize),
			nextDrop:  make(chan *container),
		},
	}, nil
}

type Server struct {
	config        *Config
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
