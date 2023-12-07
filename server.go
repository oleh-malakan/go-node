package node

import (
	"crypto/sha256"
	"net"
)

func New(address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) (*Server, error) {
	return &Server{
		nodeAddresses: nodeAddresses,
		clientCounter: newCounter(),
	}, nil
}

type Server struct {
	address       *net.UDPAddr
	nodeAddresses []*net.UDPAddr
	clientCounter *counter
	clientsLimit  int
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
	s.clientsLimit = clientsLimit

	conn, err := net.ListenUDP("udp", s.address)
	if err != nil {
		return err
	}

	go s.clientCounter.process()

	beginCore := &core{
		inData:    make(chan *incomingDatagram),
		nextDrop:  make(chan *core),
		signal:    make(chan *struct{}),
		isProcess: true,
		inProcess: s.coreBeginInProcess,
		onDestroy: coreOnDestroy,
	}
	endCore := &core{
		inData:    make(chan *incomingDatagram),
		nextDrop:  make(chan *core),
		signal:    make(chan *struct{}),
		isProcess: true,
		inProcess: coreEndInProcess,
		onDestroy: coreOnDestroy,
	}

	beginCore.next = endCore
	endCore.drop = beginCore.nextDrop

	go beginCore.process()
	go endCore.process()

	for {
		i := &incomingDatagram{
			b: make([]byte, 1432),
		}
		i.n, i.rAddr, i.err = conn.ReadFromUDP(i.b)
		if i.err != nil {

			//continue
		}
		beginCore.inData <- i
	}
}

func (s *Server) coreBeginInProcess(c *core, incoming *incomingDatagram) {
	switch {
	case incoming.b[0]&0b10000000 == 0:
		c.next.inData <- incoming
	case incoming.b[0]&0b10000000 == 1:
		if <-s.clientCounter.value <= s.clientsLimit {
			new := &core{
				inData:       make(chan *incomingDatagram),
				nextDrop:     make(chan *core),
				signal:       make(chan *struct{}),
				isProcess:    true,
				inProcess:    coreInProcess,
				onDestroy:    s.coreOnDestroy,
				next:         c.next,
				drop:         c.nextDrop,
				incoming:     incoming,
				lastIncoming: incoming,
			}
			c.next.drop = new.nextDrop
			c.next.asyncSignal()
			c.next = new
			go new.process()

			s.clientCounter.inc <- nil
		}
	}
}

func (s *Server) coreOnDestroy() {
	s.clientCounter.dec <- nil
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
