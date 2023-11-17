package node

import (
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"net"
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
	rAddr        *net.UDPAddr
	lastReadData *tReadData
	readData     *tReadData
	writeData    *tWriteData
	lastWriteMac [32]byte
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

				client := &tClient{
					conn: tls.Server(&dataport{}, tlsConfig),
					lock: make(chan *struct{}, 1),
				}

				readData.nextMac = sha256.Sum256(readData.b[1:readData.n])

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
							m[12] == readData.b[13] && m[13] == readData.b[14] && m[14] == readData.b[15] && m[15] == readData.b[16] &&
							m[16] == readData.b[17] && m[17] == readData.b[18] && m[18] == readData.b[19] && m[19] == readData.b[20] &&
							m[20] == readData.b[21] && m[21] == readData.b[22] && m[22] == readData.b[23] && m[23] == readData.b[24] &&
							m[24] == readData.b[25] && m[25] == readData.b[26] && m[26] == readData.b[27] && m[27] == readData.b[28] &&
							m[28] == readData.b[29] && m[29] == readData.b[30] && m[30] == readData.b[31] && m[31] == readData.b[32] {

							client.lastReadData.nextOk = readData.nextMac[0] == readData.b[33] && readData.nextMac[1] == readData.b[34] && readData.nextMac[2] == readData.b[35] && readData.nextMac[3] == readData.b[36] &&
								readData.nextMac[4] == readData.b[37] && readData.nextMac[5] == readData.b[38] && readData.nextMac[6] == readData.b[39] && readData.nextMac[7] == readData.b[40] &&
								readData.nextMac[8] == readData.b[41] && readData.nextMac[9] == readData.b[42] && readData.nextMac[10] == readData.b[43] && readData.nextMac[11] == readData.b[44] &&
								readData.nextMac[12] == readData.b[45] && readData.nextMac[13] == readData.b[46] && readData.nextMac[14] == readData.b[47] && readData.nextMac[15] == readData.b[48] &&
								readData.nextMac[16] == readData.b[49] && readData.nextMac[17] == readData.b[50] && readData.nextMac[18] == readData.b[51] && readData.nextMac[19] == readData.b[52] &&
								readData.nextMac[20] == readData.b[53] && readData.nextMac[21] == readData.b[54] && readData.nextMac[22] == readData.b[55] && readData.nextMac[23] == readData.b[56] &&
								readData.nextMac[24] == readData.b[57] && readData.nextMac[25] == readData.b[58] && readData.nextMac[26] == readData.b[59] && readData.nextMac[27] == readData.b[60] &&
								readData.nextMac[28] == readData.b[61] && readData.nextMac[29] == readData.b[62] && readData.nextMac[30] == readData.b[63] && readData.nextMac[31] == readData.b[64]

							client.lastReadData.next = readData
							client.lastReadData = readData

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
		}(readData)
	}
}
