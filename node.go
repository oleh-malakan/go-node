package node

import (
	"crypto/sha256"
	"crypto/tls"
	"net"
)

func Handler(nodeID string, f func(query []byte, connection *Connection)) {
	b := sha256.Sum256([]byte(nodeID))
	h := &handler{
		f: f,
	}
	_ = b[31]
	h.nodeID.p1 = uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
	h.nodeID.p2 = uint64(b[8]) | uint64(b[9])<<8 | uint64(b[10])<<16 | uint64(b[11])<<24 |
		uint64(b[12])<<32 | uint64(b[13])<<40 | uint64(b[14])<<48 | uint64(b[15])<<56
	h.nodeID.p3 = uint64(b[16]) | uint64(b[17])<<8 | uint64(b[18])<<16 | uint64(b[19])<<24 |
		uint64(b[20])<<32 | uint64(b[21])<<40 | uint64(b[22])<<48 | uint64(b[23])<<56
	h.nodeID.p4 = uint64(b[24]) | uint64(b[25])<<8 | uint64(b[26])<<16 | uint64(b[27])<<24 |
		uint64(b[28])<<32 | uint64(b[29])<<40 | uint64(b[30])<<48 | uint64(b[31])<<56

	handlers = append(handlers, h)
}

func Do(tlsConfig *tls.Config, address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) error {
	lock <- struct{}{}
	err := do(tlsConfig, handlers, address, nodeAddresses...)
	<-lock
	return err
}

type handler struct {
	nodeID tID
	f      func(query []byte, connection *Connection)
}

type tID struct {
	p1 uint64
	p2 uint64
	p3 uint64
	p4 uint64
}

var (
	handlers []*handler
	lock     chan struct{}
)

func init() {
	lock = make(chan struct{}, 1)
}

func do(tlsConfig *tls.Config, handlers []*handler, address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) error {
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		return err
	}

	type tReadData struct {
		b       []byte
		n       int
		rAddr   *net.UDPAddr
		cid     tID
		mac     tID
		nextMac tID
		err     error
	}
	type tClient struct {
		rAddr     *net.UDPAddr
		readData  []*tReadData
		readIndex int
		readed    int
	}
	var memory []*tClient

	cReadData := make(chan *tReadData, 512)
	cFreeReadData := make(chan *tReadData, 512)
	var i int
	for i = 0; i < 512; i++ {
		cFreeReadData <- &tReadData{
			b: make([]byte, 1432),
		}
	}

	go func() {
		var (
			i           int
			readData    *tReadData
			cid         tID
			client      *tClient
			clientIndex int
			iteration   int
		)
		lenMemory := len(memory)

		for {
			select {
			case readData = <-cReadData:
				if readData.err != nil {

				}

				if readData.b[0]&0x00000000 == 0 {
					client = &tClient{}
					b := sha256.Sum256(readData.b[1:readData.n])
					readData.nextMac.p1 = uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
						uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
					readData.nextMac.p2 = uint64(b[8]) | uint64(b[9])<<8 | uint64(b[10])<<16 | uint64(b[11])<<24 |
						uint64(b[12])<<32 | uint64(b[13])<<40 | uint64(b[14])<<48 | uint64(b[15])<<56
					readData.nextMac.p3 = uint64(b[16]) | uint64(b[17])<<8 | uint64(b[18])<<16 | uint64(b[19])<<24 |
						uint64(b[20])<<32 | uint64(b[21])<<40 | uint64(b[22])<<48 | uint64(b[23])<<56
					readData.nextMac.p4 = uint64(b[24]) | uint64(b[25])<<8 | uint64(b[26])<<16 | uint64(b[27])<<24 |
						uint64(b[28])<<32 | uint64(b[29])<<40 | uint64(b[30])<<48 | uint64(b[31])<<56
					client.readData = append(client.readData, readData)
				} else {

					client = nil
					clientIndex = -1
					iteration = 0
					cClientIndex := make(chan int, lenMemory)
					for i = 0; i < lenMemory; i++ {
						go func(index int) {
							j := len(memory[index].cids)
							for j > 0 {
								j--
								if memory[index].cids[j].p1 == cid.p1 && memory[index].cids[j].p2 == cid.p2 &&
									memory[index].cids[j].p3 == cid.p3 && memory[index].cids[j].p4 == cid.p4 {
									cClientIndex <- index
									return
								}
							}
							cClientIndex <- -1
						}(i)
					}

					for iteration < lenMemory {
						select {
						case i = <-cClientIndex:
							if i >= 0 {
								clientIndex = i
							}
							iteration++
						}
					}

					if clientIndex < 0 {

					}

					client = memory[clientIndex]

					cFreeReadData <- readData
				}
			}
		}
	}()

	var readData *tReadData
	for {
		readData = <-cFreeReadData
		readData.n, readData.rAddr, readData.err = conn.ReadFromUDP(readData.b)
		cReadData <- readData
	}
}
