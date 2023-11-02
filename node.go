package node

import (
	"crypto/sha256"
	"net"
)

func Handler(nodeID string, f func(query []byte, connection *Connection)) {
	b := sha256.Sum256([]byte(nodeID))
	h := &handler{
		f: f,
	}
	_ = b[31]
	h.nodeID[0] = uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
	h.nodeID[1] = uint64(b[8]) | uint64(b[9])<<8 | uint64(b[10])<<16 | uint64(b[11])<<24 |
		uint64(b[12])<<32 | uint64(b[13])<<40 | uint64(b[14])<<48 | uint64(b[15])<<56
	h.nodeID[2] = uint64(b[16]) | uint64(b[17])<<8 | uint64(b[18])<<16 | uint64(b[19])<<24 |
		uint64(b[20])<<32 | uint64(b[21])<<40 | uint64(b[22])<<48 | uint64(b[23])<<56
	h.nodeID[3] = uint64(b[24]) | uint64(b[25])<<8 | uint64(b[26])<<16 | uint64(b[27])<<24 |
		uint64(b[28])<<32 | uint64(b[29])<<40 | uint64(b[30])<<48 | uint64(b[31])<<56

	handlers = append(handlers, h)
}

func Do(address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) error {
	lock <- struct{}{}
	err := do(handlers, address, nodeAddresses...)
	<-lock
	return err
}

type handler struct {
	nodeID [4]uint64
	f      func(query []byte, connection *Connection)
}

var (
	handlers []*handler
	lock     chan struct{}
)

func init() {
	lock = make(chan struct{}, 1)
}

func do(handlers []*handler, address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) error {

	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		return err
	}

	type tClient struct {
		rAddr *net.UDPAddr
		cids  [][4]uint64
	}
	type tMemory struct {
		clients   []*tClient
		index     []int
		freeIndex []int
	}
	memory := tMemory{}
	tmpMemory := tMemory{}

	type tReadData struct {
		b     []byte
		n     int
		rAddr *net.UDPAddr
		err   error
	}

	cReadData := make(chan *tReadData, 512)
	cFreeReadData := make(chan *tReadData, 512)
	for i := 0; i < 512; i++ {
		cFreeReadData <- &tReadData{
			b: make([]byte, 560),
		}
	}

	go func() {
		var cli *tClient
		for i := 0; i < len(memory.clients); i++ {
			cli = memory.clients[memory.index[i]]
			if cli.rAddr.IP.Equal(rAddr.IP) && cli.rAddr.Port == rAddr.Port {
				break
			}
			cli = nil
		}
		if cli == nil {

		}
	}()

	var readData *tReadData
	for {
		readData = <-cFreeReadData
		readData.n, readData.rAddr, readData.err = conn.ReadFromUDP(readData.b)
		cReadData <- readData
	}
}
