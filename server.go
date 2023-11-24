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
		address:       address,
		nodeAddresses: nodeAddresses,
		in:            make(chan *incomingPackage, 512),
		nextDrop:      make(chan *container),
	}, nil
}

type Server struct {
	tlsConfig     *tls.Config
	address       *net.UDPAddr
	nodeAddresses []*net.UDPAddr

	next     *container
	in       chan *incomingPackage
	nextDrop chan *container
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
		incoming := &incomingPackage{
			b: make([]byte, 1432),
		}
		incoming.n, incoming.rAddr, incoming.err = conn.ReadFromUDP(incoming.b)

		s.in <- incoming
	}
}

func (s *Server) process() {
	for {
		select {
		case p := <-s.in:
			switch {
			case p.b[0]>>7&1 == 0:
				p.nextMac = sha256.Sum256(p.b[1:p.n])
				if s.next != nil {
					s.next.in <- p
				} else {
					s.next = newContainer(p, s.nextDrop, s.tlsConfig)
					go s.next.process()
				}
			case p.b[0]>>7&1 == 1:
				p.nextMac = sha256.Sum256(p.b[65:p.n])
				if s.next != nil {
					s.next.in <- p
				}
			}
		case dropNode := <-s.nextDrop:
			s.next = dropNode.next
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
