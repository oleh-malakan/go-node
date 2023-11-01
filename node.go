package node

import (
	"crypto/sha256"
	"encoding/binary"
	"net"
)

func Handler(nodeID string, f func(query []byte, connection *Connection)) {
	b := sha256.Sum256([]byte(nodeID))
	handler := &handler{
		f:      f,
	}
	for i := 0; i < 4; i++ {
		handler.nodeID[i] = binary.LittleEndian.Uint64(b[i * 8: i * 8 + 8])
	}

	initHandlers = append(initHandlers, handler)
}

func Do(address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) error {
	handlers = make([]*handler, len(initHandlers))
	copy(handlers, initHandlers)
	initHandlers = nil

	return nil
}

type handler struct {
	nodeID [4]uint64
	f      func(query []byte, connection *Connection)
}

var (
	initHandlers []*handler
	handlers     []*handler
)
