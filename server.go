package node

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"net"
)

func New(address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) (*Server, error) {
	return &Server{
		nodeAddresses: nodeAddresses,
		transport:     &transport{},
	}, nil
}

type Server struct {
	address        *net.UDPAddr
	nodeAddresses  []*net.UDPAddr
	transport      *transport
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
func (s *Server) Run(connectionsLimit int) error {
	if connectionsLimit <= 0 {
		connectionsLimit = 524288
	}

	var err error
	s.transport.conn, err = net.ListenUDP("udp", s.address)
	if err != nil {
		return err
	}

	s.privateKey, err = ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	s.publicKeyBytes = s.privateKey.PublicKey().Bytes()

	in := make(chan *datagram)
	go s.process(in, connectionsLimit)

	for {
		i := &datagram{
			b: make([]byte, 1432),
		}
		i.n, i.rAddr, i.err = s.transport.read(i.b)
		if i.err != nil {

			//continue
		}
		in <- i
	}
}

func (s *Server) process(in chan *datagram, connectionsLimit int) {
	controller := &controller{
		connectionsLimit: connectionsLimit,
		counter:          newCounter(),
		drop:             make(chan int),
	}
	go controller.counter.process()
	for {
		select {
		case i := <-in:
			if i.b[0] != 0 {
				controller.in(i)
			} else {
				s.serverHello(i, controller.cid())
			}
		case i := <-controller.drop:
			controller.free(i)
		}
	}
}

func (s *Server) serverHello(incoming *datagram, cid []byte) {
	b := make([]byte, datagramMinLen)
	copy(b[1:33], s.publicKeyBytes)

	rKey, err := ecdh.P256().NewPublicKey(incoming.b[1:33])
	if err != nil {
		return
	}

	secret, err := s.privateKey.ECDH(rKey)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	rand.Reader.Read(b[81:93])
	aead.Seal(b[:33], b[81:93], cid, b[:33])
	// error b serverHello
	_, err = s.transport.write(b, incoming.rAddr)
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
