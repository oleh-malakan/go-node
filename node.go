package node

import (
	"crypto/sha256"
	"crypto/tls"
	"net"
)

type Config struct {
	Strokes     int
	ClientCount int
	BufferSize  int
}

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

func Do(tlsConfig *tls.Config, config *Config, address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) error {
	if tlsConfig == nil {
		return newError("require tls config")
	}

	if config == nil {
		config = &Config{
			Strokes:     4,
			ClientCount: 9096,
			BufferSize:  9096,
		}
	}

	return do(handlers, tlsConfig, config.Strokes, config.ClientCount, config.BufferSize, address, nodeAddresses...)
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
)

func do(handlers []*handler, tlsConfig *tls.Config,
	strokes int, clientCount int, bufferSize int,
	address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) error {
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		return err
	}

	cStroke := make(chan *struct{}, strokes)
	go func() {
		for i := 0; i < strokes; i++ {
			cStroke <- nil
		}
	}()

	type tReadData struct {
		b       []byte
		n       int
		rAddr   *net.UDPAddr
		mac     tID
		nextMac tID
		next    *tReadData
		err     error
	}

	type tWriteData struct {
		prevMac tID
		mac     tID
		prev    *tWriteData
	}

	cFreeReadData := make(chan *tReadData, bufferSize)
	go func() {
		for i := 0; i < bufferSize; i++ {
			cFreeReadData <- &tReadData{
				b: make([]byte, 1432),
			}
		}
	}()

	type tClient struct {
		conn         *tls.Conn
		rAddr        *net.UDPAddr
		lastReadData *tReadData
		readData     *tReadData
		nextReadMac  tID
		readed       int
		writeData    *tWriteData
		lastWriteMac tID
		lock         chan *struct{}
		bypassLock   chan *struct{}
		next         *tClient
		drop         bool
	}

	type tReportIterationCount struct {
		count int
	}

	var memory *tClient
	memoryLock := make(chan *struct{}, 1)

	reportFoundClient := new(struct{})

	cFreeClient := make(chan *tClient, clientCount)
	go func() {
		for i := 0; i < clientCount; i++ {
			cFreeClient <- &tClient{
				conn:       tls.Server(&dataport{}, tlsConfig),
				lock:       make(chan *struct{}, 1),
				bypassLock: make(chan *struct{}, 1),
			}
		}
	}()

	bypassMemory := func(readData *tReadData) {
		var (
			cid      tID
			client   *tClient
			bNextMac [32]byte
		)

		switch {
		case readData.b[0]&0b00000000 == 0b00000000:
			select {
			case client = <-cFreeClient:
				bNextMac = sha256.Sum256(readData.b[1:readData.n])
				readData.nextMac.p1 = uint64(bNextMac[0]) | uint64(bNextMac[1])<<8 | uint64(bNextMac[2])<<16 | uint64(bNextMac[3])<<24 |
					uint64(bNextMac[4])<<32 | uint64(bNextMac[5])<<40 | uint64(bNextMac[6])<<48 | uint64(bNextMac[7])<<56
				readData.nextMac.p2 = uint64(bNextMac[8]) | uint64(bNextMac[9])<<8 | uint64(bNextMac[10])<<16 | uint64(bNextMac[11])<<24 |
					uint64(bNextMac[12])<<32 | uint64(bNextMac[13])<<40 | uint64(bNextMac[14])<<48 | uint64(bNextMac[15])<<56
				readData.nextMac.p3 = uint64(bNextMac[16]) | uint64(bNextMac[17])<<8 | uint64(bNextMac[18])<<16 | uint64(bNextMac[19])<<24 |
					uint64(bNextMac[20])<<32 | uint64(bNextMac[21])<<40 | uint64(bNextMac[22])<<48 | uint64(bNextMac[23])<<56
				readData.nextMac.p4 = uint64(bNextMac[24]) | uint64(bNextMac[25])<<8 | uint64(bNextMac[26])<<16 | uint64(bNextMac[27])<<24 |
					uint64(bNextMac[28])<<32 | uint64(bNextMac[29])<<40 | uint64(bNextMac[30])<<48 | uint64(bNextMac[31])<<56

				client.readData = readData
				client.nextReadMac = readData.nextMac
				client.drop = false
				go func() {
					if err := client.conn.Handshake(); err != nil {

					}
				}()

				memoryLock <- nil
				if memory != nil {
					client.next = memory
					memory = client
				} else {
					memory = client
				}
				<-memoryLock
			default:
				cFreeReadData <- readData
			}
		case readData.b[0]&0b10000000 == 0b10000000 && memory != nil:
			cid.p1 = uint64(readData.b[1]) | uint64(readData.b[2])<<8 | uint64(readData.b[3])<<16 | uint64(readData.b[4])<<24 |
				uint64(readData.b[5])<<32 | uint64(readData.b[6])<<40 | uint64(readData.b[7])<<48 | uint64(readData.b[8])<<56
			cid.p2 = uint64(readData.b[9]) | uint64(readData.b[10])<<8 | uint64(readData.b[11])<<16 | uint64(readData.b[12])<<24 |
				uint64(readData.b[13])<<32 | uint64(readData.b[14])<<40 | uint64(readData.b[15])<<48 | uint64(readData.b[16])<<56
			cid.p3 = uint64(readData.b[17]) | uint64(readData.b[18])<<8 | uint64(readData.b[19])<<16 | uint64(readData.b[20])<<24 |
				uint64(readData.b[21])<<32 | uint64(readData.b[22])<<40 | uint64(readData.b[23])<<48 | uint64(readData.b[24])<<56
			cid.p4 = uint64(readData.b[25]) | uint64(readData.b[26])<<8 | uint64(readData.b[27])<<16 | uint64(readData.b[28])<<24 |
				uint64(readData.b[29])<<32 | uint64(readData.b[30])<<40 | uint64(readData.b[31])<<48 | uint64(readData.b[32])<<56
			readData.mac.p1 = uint64(readData.b[33]) | uint64(readData.b[34])<<8 | uint64(readData.b[35])<<16 | uint64(readData.b[36])<<24 |
				uint64(readData.b[37])<<32 | uint64(readData.b[38])<<40 | uint64(readData.b[39])<<48 | uint64(readData.b[40])<<56
			readData.mac.p2 = uint64(readData.b[41]) | uint64(readData.b[42])<<8 | uint64(readData.b[43])<<16 | uint64(readData.b[44])<<24 |
				uint64(readData.b[45])<<32 | uint64(readData.b[46])<<40 | uint64(readData.b[47])<<48 | uint64(readData.b[48])<<56
			readData.mac.p3 = uint64(readData.b[49]) | uint64(readData.b[50])<<8 | uint64(readData.b[51])<<16 | uint64(readData.b[52])<<24 |
				uint64(readData.b[53])<<32 | uint64(readData.b[54])<<40 | uint64(readData.b[55])<<48 | uint64(readData.b[56])<<56
			readData.mac.p4 = uint64(readData.b[57]) | uint64(readData.b[58])<<8 | uint64(readData.b[59])<<16 | uint64(readData.b[60])<<24 |
				uint64(readData.b[61])<<32 | uint64(readData.b[62])<<40 | uint64(readData.b[63])<<48 | uint64(readData.b[64])<<56
			bNextMac = sha256.Sum256(readData.b[65:readData.n])
			readData.nextMac.p1 = uint64(bNextMac[0]) | uint64(bNextMac[1])<<8 | uint64(bNextMac[2])<<16 | uint64(bNextMac[3])<<24 |
				uint64(bNextMac[4])<<32 | uint64(bNextMac[5])<<40 | uint64(bNextMac[6])<<48 | uint64(bNextMac[7])<<56
			readData.nextMac.p2 = uint64(bNextMac[8]) | uint64(bNextMac[9])<<8 | uint64(bNextMac[10])<<16 | uint64(bNextMac[11])<<24 |
				uint64(bNextMac[12])<<32 | uint64(bNextMac[13])<<40 | uint64(bNextMac[14])<<48 | uint64(bNextMac[15])<<56
			readData.nextMac.p3 = uint64(bNextMac[16]) | uint64(bNextMac[17])<<8 | uint64(bNextMac[18])<<16 | uint64(bNextMac[19])<<24 |
				uint64(bNextMac[20])<<32 | uint64(bNextMac[21])<<40 | uint64(bNextMac[22])<<48 | uint64(bNextMac[23])<<56
			readData.nextMac.p4 = uint64(bNextMac[24]) | uint64(bNextMac[25])<<8 | uint64(bNextMac[26])<<16 | uint64(bNextMac[27])<<24 |
				uint64(bNextMac[28])<<32 | uint64(bNextMac[29])<<40 | uint64(bNextMac[30])<<48 | uint64(bNextMac[31])<<56
		default:
			cFreeReadData <- readData
		}

		memory := memory
		if memory != nil {
			cReport := make(chan interface{})
			reportIterationCount := &tReportIterationCount{}

			var bypass func(c *tClient, i int)
			bypass = func(c *tClient, i int) {
				i++
				c.bypassLock <- nil
			NEXT:
				if c.next != nil {
					if !c.next.drop {
						go bypass(c.next, i)
					} else {
						d := c.next
						d.bypassLock <- nil
						d.lock <- nil
						if d.next != nil {
							c.next = d.next
							d.next = nil
						} else {
							c.next = nil
						}
						cFreeClient <- d
						<-d.lock
						<-d.bypassLock
						goto NEXT
					}
				} else {
					reportIterationCount.count = i
					cReport <- reportIterationCount
				}
				<-c.bypassLock

				c.lock <- nil
				if !c.drop {
					if readData != nil {
						w := c.writeData
						m := c.lastWriteMac
					LOOP:
						if m.p1 == cid.p1 && m.p2 == cid.p2 &&
							m.p3 == cid.p3 && m.p4 == cid.p4 {

							//
							//
							//

							cReport <- reportFoundClient
							goto FOUND
						}
						if w != nil && w.prev != nil {
							w = w.prev
							m = w.mac

							goto LOOP
						}
					FOUND:
					}
				}
				<-c.lock
				cReport <- nil
			}

			go bypass(memory, 0)

			var (
				report         interface{}
				iteration      int
				iterationCount int
				foundClient    bool
			)

			for iterationCount == 0 || iteration < iterationCount {
				select {
				case report = <-cReport:
					switch report {
					case nil:
						iteration++
					case reportFoundClient:
						foundClient = true
					case reportIterationCount:
						iterationCount = reportIterationCount.count
					}
				}
			}

			if !foundClient && readData != nil {
				cFreeReadData <- readData
			}
		}
	}

	var readData *tReadData
	for {
		select {
		case <-cStroke:
			go func() {
				readData = <-cFreeReadData
				readData.n, readData.rAddr, readData.err = conn.ReadFromUDP(readData.b)
				go bypassMemory(readData)
				<-cStroke
			}()
		}
	}
}
