package node

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"net"

	"github.com/oleh-malakan/go-node/memory"
)

// connectionsLimit default value 2000000 if 0, max value 1000000000
func Run(connectionsLimit int, address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) (*Server, error) {
	if connectionsLimit <= 0 {
		connectionsLimit = 2000000
	}
	if connectionsLimit > 1000000000 {
		connectionsLimit = 1000000000
	}

	s := &Server{
		address:       address,
		nodeAddresses: nodeAddresses,
	}
	var err error
	s.transport.conn, err = net.ListenUDP("udp", s.address)
	if err != nil {
		return nil, err
	}

	go s.process()

	return s, nil
}

type Server struct {
	address          *net.UDPAddr
	nodeAddresses    []*net.UDPAddr
	transport        *transport
	connectionsLimit int
}

func (s *Server) Listen(nodeID string) (*Listener, error) {
	return &Listener{
		nodeID: sha256.Sum224([]byte(nodeID)),
	}, nil
}

func (s *Server) process() {
	in := make(chan *datagram)
	go func() {
		for {
			i := &datagram{
				b: make([]byte, 1432),
			}
			i.n, i.rAddr, i.err = s.transport.readUDP(i.b)
			if i.err != nil {

				//continue
			}
			in <- i
		}
	}()

	memory := &memory.Memory[core]{}
	drop := make(chan *core)
	for {
		select {
		case i := <-in:
			if i.b[0] != 0 {
				cIDDatagram := parseCIDDatagram(i)
				if current := memory.Get(cIDDatagram.cid); current != nil {
					current.inData <- cIDDatagram.datagram
				}
			} else {
				if len(i.b) >= datagramMinLen && memory.Len() < s.connectionsLimit {
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
					new.cipher, err = cipher.NewGCM(block)
					if err != nil {
						continue
					}

					new.cid = memory.Put(new)

					b := make([]byte, datagramMinLen)
					copy(b[1:33], privateKey.PublicKey().Bytes())
					rand.Reader.Read(b[53:65])

					binary.BigEndian.PutUint32(b[33:], uint32(new.cid))

					new.cipher.Seal(b[:33], b[55:67], b[33:37], b[:33])
					_, err = s.transport.writeUDP(b, i.rAddr)
					if err != nil {
						memory.Free(new.cid)
						continue
					}

					go new.process()
				}
			}
		case core := <-drop:
			memory.Free(core.cid)
		}
	}
}

func (s *Server) Close() error {
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
