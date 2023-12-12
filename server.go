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
		transport:     &transport{},
	}, nil
}

type Server struct {
	address       *net.UDPAddr
	nodeAddresses []*net.UDPAddr
	transport     *transport
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
	var connectionCount int
	memory := &indexArray[core]{}
	drop := make(chan int64)
	for {
		select {
		case i := <-in:
			if i.b[0] != 0 {
				cIDDatagram := parseCIDDatagram(i)
				if current := memory.get(cIDDatagram.cid); current != nil {
					current.inData <- cIDDatagram.datagram
				}
			} else {
				if connectionCount < connectionsLimit && len(i.b) >= datagramMinLen {
					privateKey, err := ecdh.X25519().GenerateKey(rand.Reader)
					if err != nil {
						continue
					}
					remotePublicKey, err := ecdh.X25519().NewPublicKey(i.b[1:33])
					if err != nil {
						continue
					}
					secret, err := privateKey.ECDH(remotePublicKey)
					if err != nil {
						continue
					}
					block, err := aes.NewCipher(secret)
					if err != nil {
						continue
					}

					i.did = int64(binary.BigEndian.Uint64(i.b[25:33]))
					i.begin = i.n
					new := &core{
						inData:       make(chan *datagram),
						drop:         drop,
						isProcess:    true,
						incoming:     i,
						lastIncoming: i,
					}
					new.aead, err = cipher.NewGCM(block)
					if err != nil {
						continue
					}
					new.cid = memory.put(new)

					b := make([]byte, datagramMinLen)
					copy(b[1:33], privateKey.PublicKey().Bytes())
					rand.Reader.Read(b[57:69])
					var cidB []byte
					cidB = binary.BigEndian.AppendUint64(cidB, uint64(new.cid))

					new.aead.Seal(b[:33], b[57:69], cidB, b[:33])
					_, err = s.transport.write(b, i.rAddr)
					if err != nil {
						memory.free(new.cid)
						continue
					}

					go new.process()
					connectionCount++
				}
			}
		case i := <-drop:
			memory.free(i)
			connectionCount--
		}
	}
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
