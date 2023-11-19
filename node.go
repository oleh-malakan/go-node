package node

import (
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"net"
	"time"
)

func Handler(nodeID string, f func(query []byte, connection *Connection)) {
	h := &handler{
		f: f,
	}
	h.nodeID = sha256.Sum256([]byte(nodeID))

	handlers = append(handlers, h)
}

func Do(tlsConfig *tls.Config, address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) error {
	if tlsConfig == nil {
		return errors.New("require tls config")
	}

	return do(handlers, tlsConfig, address, nodeAddresses...)
}

type handler struct {
	nodeID [32]byte
	f      func(query []byte, connection *Connection)
}

var (
	handlers []*handler
)

type tReadData struct {
	b      []byte
	n      int
	readed int
	rAddr  *net.UDPAddr
	//	mac     [32]byte
	nextMac [32]byte
	next    *tReadData
	nextOk  bool
	err     error
}

type tWriteData struct {
	//prevMac [32]byte
	mac  [32]byte
	prev *tWriteData
}

type tClient struct {
	conn         *tls.Conn
	lastReadData *tReadData
	readData     *tReadData
	heap         *tHeap
	writeData    *tWriteData
	next         *tClient
	lock         chan *struct{}
	drop         bool
}

func do(handlers []*handler, tlsConfig *tls.Config, address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) error {
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		return err
	}

	var memory *tClient
	memoryLock := make(chan *struct{}, 1)

	for {
		readData := &tReadData{
			b: make([]byte, 1432),
		}
		readData.n, readData.rAddr, readData.err = conn.ReadFromUDP(readData.b)

		go func(readData *tReadData) {
			switch {
			case readData.b[0]>>7&1 == 0:
				readData.nextMac = sha256.Sum256(readData.b[1:readData.n])
				client := &tClient{
					conn: tls.Server(&dataport{}, tlsConfig),
					lock: make(chan *struct{}, 1),
					heap: &tHeap{
						cap:     512,
						timeout: int64(time.Duration(50) * time.Millisecond),
					},
				}
				client.readData = readData
				client.lastReadData = readData
				client.drop = false

				memoryLock <- nil
				if memory != nil {
					client.next = memory
					memory = client
				} else {
					memory = client
				}
				<-memoryLock
			case readData.b[0]>>7&1 == 1 && memory != nil:
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
						for w != nil {
							if compareID(w.mac[0:32], readData.b[1:33]) {
								if client.lastReadData.nextOk = compareID(client.lastReadData.nextMac[0:32], readData.b[33:65]); client.lastReadData.nextOk {
									client.lastReadData.next = readData
									client.lastReadData = readData
									client.lastReadData.next = client.heap.find(readData.nextMac[0:32])
									client.lastReadData.nextOk = client.lastReadData.next != nil
								} else {
									client.heap.put(readData)
								}

								next = nil
								goto FOUND
							}
							w = w.prev
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
		}(readData)
	}
}
