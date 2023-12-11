package node

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"net"

	"github.com/oleh-malakan/go-node/internal"
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
	privateKey     *ecdh.PrivateKey
	publicKeyBytes []byte
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

	inData := make(chan *datagram)
	go s.process(inData, connectionsLimit)

	for {
		i := &datagram{
			b: make([]byte, 1432),
		}
		i.n, i.rAddr, i.err = s.transport.read(i.b)
		if i.err != nil {

			//continue
		}
		inData <- i
	}
}

func (s *Server) process(inData chan *datagram, connectionsLimit int) {
	var connectionsCount int
	memory := &internal.IndexArray[core]{}
	cidManager := &internal.CIDManager{}
	drop := make(chan int)
	for {
		select {
		case i := <-inData:
			if i.b[0] != 0 {
				cIDDatagram := parseCIDDatagram(i)
				if c := memory.Get(cIDDatagram.index); c != nil && c.checkCID(cIDDatagram.cid) {
					c.inData <- i
				} else {
					if connectionsCount < connectionsLimit && cidManager.Put(cIDDatagram.cid) {
						new := &core{
							inData:       make(chan *datagram),
							drop:         drop,
							isProcess:    true,
							incoming:     i,
							lastIncoming: i,
						}
						new.index = memory.Put(new)
						go new.process()
						connectionsCount++
					}
				}
			} else {
				s.serverHello(i, cidManager.CID())
			}
		case i := <-drop:
			memory.Free(i)
			connectionsCount--
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
