package node

import (
	"crypto/sha256"
	"crypto/tls"
	"net"
)

func Handler(nodeID string, f func(query []byte, connection *Connection)) {
	h := &handler{
		f: f,
	}
	h.nodeID = sha256.Sum256([]byte(nodeID))

	handlers = append(handlers, h)
}

func Do(tlsConfig *tls.Config, threads int, address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) error {
	if tlsConfig == nil {
		return newError("require tls config")
	}

	return do(handlers, tlsConfig, threads, address, nodeAddresses...)
}

type handler struct {
	nodeID [32]byte
	f      func(query []byte, connection *Connection)
}

var (
	handlers []*handler
)

type tReadData struct {
	b     []byte
	n     int
	rAddr *net.UDPAddr
	//	mac     [32]byte
	nextMac [32]byte
	next    *tReadData
	err     error
}

type tWriteData struct {
	//prevMac [32]byte
	mac  [32]byte
	prev *tWriteData
}

type tClient struct {
	conn         *tls.Conn
	rAddr        *net.UDPAddr
	lastReadData *tReadData
	readData     *tReadData
	nextReadMac  [32]byte
	readed       int
	writeData    *tWriteData
	lastWriteMac [32]byte
	next         *tClient
	lock         chan *struct{}
	drop         bool
}

func do(handlers []*handler, tlsConfig *tls.Config, threads int,
	address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) error {
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		return err
	}

	cFatal := make(chan error)

	var memory *tClient
	memoryLock := make(chan *struct{}, 1)

	for i := 0; i < threads; i++ {
		go func() {
			for {
				readData := &tReadData{
					b: make([]byte, 1432),
				}
				readData.n, readData.rAddr, readData.err = conn.ReadFromUDP(readData.b)

				switch {
				case readData.b[0]&0b00000000 == 0b00000000:
					client := &tClient{
						conn: tls.Server(&dataport{}, tlsConfig),
						lock: make(chan *struct{}, 1),
					}
					readData.nextMac = sha256.Sum256(readData.b[1:readData.n])

					client.readData = readData
					client.nextReadMac = readData.nextMac
					client.drop = false
					/*
						go func() {
							if err := client.conn.Handshake(); err != nil {

							}
						}()
					*/

					memoryLock <- nil
					if memory != nil {
						client.next = memory
						memory = client
					} else {
						memory = client
					}
					<-memoryLock
				case readData.b[0]&0b10000000 == 0b10000000 && memory != nil:
					/*
						readData.mac = readData.b[33:65]
					*/
					readData.nextMac = sha256.Sum256(readData.b[65:readData.n])
					var client *tClient
					memoryLock <- nil
					if memory != nil {
						client = memory
						client.lock <- nil
					}
					<-memoryLock
					for client != nil {
						var next *tClient

						if !client.drop {
							w := client.writeData
							m := client.lastWriteMac
						LOOP:
							if m[0] == readData.b[1] && m[1] == readData.b[2] && m[2] == readData.b[3] && m[3] == readData.b[4] &&
								m[4] == readData.b[5] && m[5] == readData.b[6] && m[6] == readData.b[7] && m[7] == readData.b[8] &&
								m[8] == readData.b[9] && m[9] == readData.b[10] && m[10] == readData.b[11] && m[11] == readData.b[12] &&
								m[12] == readData.b[13] && m[13] == readData.b[14] && m[14] == readData.b[15] && m[15] == readData.b[16] {

								//
								//
								//

								goto FOUND
							}
							if w != nil && w.prev != nil {
								w = w.prev
								m = w.mac

								goto LOOP
							}
						}

						next = client.next
						if next != nil {
							next.lock <- nil
						}
					FOUND:
						<-client.lock
						client = next
					}
				}
			}
		}()
	}

	return <-cFatal
}
