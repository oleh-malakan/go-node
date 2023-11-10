package node

import (
	"crypto/sha256"
	"crypto/tls"
	"net"
)

type Config struct {
	ClientsLimit int
	BufferSize   int
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
			ClientsLimit: 0,
			BufferSize:   512,
		}
	}

	return do(handlers, tlsConfig, config.ClientsLimit, config.BufferSize, address, nodeAddresses...)
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
	clientsLimit int, bufferSize int,
	address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) error {
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		return err
	}

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

	cReadData := make(chan *tReadData, bufferSize)
	cFreeReadData := make(chan *tReadData, bufferSize)
	for i := 0; i < bufferSize; i++ {
		cFreeReadData <- &tReadData{
			b: make([]byte, 1432),
		}
	}

	go func() {
		type tClient struct {
			rAddr        *net.UDPAddr
			lastReadData *tReadData
			readData     *tReadData
			nextReadMac  tID
			readed       int
			writeData    *tWriteData
			lastWriteMac tID
			next         *tClient
			drop         bool
		}

		var (
			memory            *tClient
			lenMemory         int
			readData          *tReadData
			cid               tID
			client            *tClient
			bypass            func(c *tClient)
			cBypass           chan *struct{}
			bypassDone        *struct{}
			bypassFoundClient *struct{}
			bypassDropClient  *struct{}
			foundClient       bool
			dropClientCount   int
			iteration         int
			bNextMac          [32]byte
		)

		bypassFoundClient = new(struct{})
		bypass = func(c *tClient) {
		NEXT:
			if c.next != nil && !c.next.drop {
				go bypass(c.next)
			} else if c.next.next != nil {
				c.next = c.next.next
				goto NEXT
			} else {
				c.next = nil
			}

			if readData != nil {
				w := c.writeData
				m := c.lastWriteMac
			LOOP:
				if m.p1 == cid.p1 && m.p2 == cid.p2 &&
					m.p3 == cid.p3 && m.p4 == cid.p4 {

					//
					//
					//
					readData = nil

					cBypass <- bypassFoundClient

					return
				}
				if w != nil && w.prev != nil {
					w = w.prev
					m = w.mac

					goto LOOP
				}
			}

			cBypass <- nil
		}

		goto LOOP
	BYPASS:
		go bypass(memory)

		iteration = 0
		foundClient = false
		for iteration < lenMemory {
			select {
			case bypassDone = <-cBypass:
				switch bypassDone {
				case bypassFoundClient:
					foundClient = true
				case bypassDropClient:
					dropClientCount++
				}
				iteration++
			}
		}

		if !foundClient && readData != nil {
			cFreeReadData <- readData
			readData = nil
		}

		if dropClientCount > 0 {
			lenMemory = lenMemory - dropClientCount
			cBypass = make(chan *struct{}, lenMemory)
			dropClientCount = 0
		}
	LOOP:
		for {
			select {
			case readData = <-cReadData:
				if readData.err != nil {
					cFreeReadData <- readData
					readData = nil
				}

				switch {
				case readData.b[0]&0b00000000 == 0b00000000:
					if clientsLimit > 0 {
						if lenMemory >= clientsLimit {
							// go f(readData.rAddr)
							cFreeReadData <- readData
							readData = nil
							break
						}
					}

					bNextMac = sha256.Sum256(readData.b[1:readData.n])
					readData.nextMac.p1 = uint64(bNextMac[0]) | uint64(bNextMac[1])<<8 | uint64(bNextMac[2])<<16 | uint64(bNextMac[3])<<24 |
						uint64(bNextMac[4])<<32 | uint64(bNextMac[5])<<40 | uint64(bNextMac[6])<<48 | uint64(bNextMac[7])<<56
					readData.nextMac.p2 = uint64(bNextMac[8]) | uint64(bNextMac[9])<<8 | uint64(bNextMac[10])<<16 | uint64(bNextMac[11])<<24 |
						uint64(bNextMac[12])<<32 | uint64(bNextMac[13])<<40 | uint64(bNextMac[14])<<48 | uint64(bNextMac[15])<<56
					readData.nextMac.p3 = uint64(bNextMac[16]) | uint64(bNextMac[17])<<8 | uint64(bNextMac[18])<<16 | uint64(bNextMac[19])<<24 |
						uint64(bNextMac[20])<<32 | uint64(bNextMac[21])<<40 | uint64(bNextMac[22])<<48 | uint64(bNextMac[23])<<56
					readData.nextMac.p4 = uint64(bNextMac[24]) | uint64(bNextMac[25])<<8 | uint64(bNextMac[26])<<16 | uint64(bNextMac[27])<<24 |
						uint64(bNextMac[28])<<32 | uint64(bNextMac[29])<<40 | uint64(bNextMac[30])<<48 | uint64(bNextMac[31])<<56

					client = &tClient{
						readData:    readData,
						nextReadMac: readData.nextMac,
					}
					readData = nil

					if memory != nil {
						client.next = memory
						memory = client
					} else {
						memory = client
					}
					lenMemory++
					cBypass = make(chan *struct{}, lenMemory)
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

					goto BYPASS
				default:
					cFreeReadData <- readData
					readData = nil
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
