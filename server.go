package node

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"net"
)

func New(address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) (*Server, error) {
	return &Server{
		nodeAddresses: nodeAddresses,
		clientCounter: newCounter(),
		transport:     &transport{},
	}, nil
}

type Server struct {
	address        *net.UDPAddr
	nodeAddresses  []*net.UDPAddr
	privateKey     *ecdh.PrivateKey
	publicKeyBytes []byte
	transport      *transport
	clientCounter  *counter
	clientsLimit   int
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

	go s.clientCounter.process()

	beginCore := &core{
		inData:         make(chan *incomingDatagram),
		drop:           make(chan *core, 1),
		isProcess:      true,
		inProcess:      s.coreBeginInProcess,
		destroyProcess: coreDestroyProcess,
	}
	endCore := &core{
		inData:         make(chan *incomingDatagram),
		drop:           make(chan *core),
		isProcess:      true,
		inProcess:      s.coreEndInProcess,
		destroyProcess: coreDestroyProcess,
	}
	beginCore.next = endCore
	endCore.next = beginCore

	go beginCore.process()
	go endCore.process()

	for {
		i := &incomingDatagram{
			cipherB: make([]byte, 1432),
		}
		i.n, i.rAddr, i.err = s.transport.read(i.cipherB)
		if i.err != nil {

			//continue
		}
		beginCore.inData <- i
	}
}

func (s *Server) newClientID() []byte {
	var ID []byte

	return ID[:]
}

func (s *Server) serverHello(incoming *incomingDatagram) {
	b := make([]byte, datagramMinLen)
	copy(b[1:33], s.publicKeyBytes)

	rKey, err := ecdh.P256().NewPublicKey(incoming.cipherB[1:33])
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

	binary.BigEndian.PutUint64(b[33:41], uint64(aead.NonceSize()))
	openDataLen := 41 + aead.NonceSize()
	rand.Reader.Read(b[41:openDataLen])
	aead.Seal(b[:openDataLen], b[41:openDataLen], s.newClientID(), b[:openDataLen])
	_, err = s.transport.write(b, incoming.rAddr)
}

func (s *Server) decodeClientHello2(incoming *incomingDatagram) bool {
	return true
}

func (s *Server) coreBeginInProcess(c *core, incoming *incomingDatagram) {
	if incoming.cipherB[0]&0b10000000 != 0 {
		if incoming.b == nil {
			incoming.prepareCID()
			c.next.inData <- incoming
		} else {
			if <-s.clientCounter.value <= s.clientsLimit {
				new := &core{
					inData:         make(chan *incomingDatagram),
					drop:           make(chan *core),
					isProcess:      true,
					inProcess:      s.coreInProcess,
					destroyProcess: s.coreDestroyProcess,
					next:           c.next,
					incoming:       incoming,
					lastIncoming:   incoming,
				}
				c.next = new
				go new.process()
				s.clientCounter.inc <- nil
			}
		}
	} else {
		s.serverHello(incoming)
	}
}

func (s *Server) coreInProcess(core *core, incoming *incomingDatagram) {
	if core.cid.ID1 != incoming.cid.ID1 || core.cid.ID2 != incoming.cid.ID2 ||
		core.cid.ID3 != incoming.cid.ID3 || core.cid.ID4 != incoming.cid.ID4 {
		core.next.inData <- incoming
	} else if incoming.cipherB[0]&0b01000000 == 0 {
		coreInProcess(core, incoming)
	}
}

func (s *Server) coreEndInProcess(core *core, incoming *incomingDatagram) {
	if incoming.cipherB[0]&0b01000000 != 0 {
		if s.decodeClientHello2(incoming) {
			core.next.inData <- incoming
		}
	}
}

func (s *Server) coreDestroyProcess() {
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
