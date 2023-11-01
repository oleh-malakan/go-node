package node

import (
	"crypto/sha256"
	"encoding/binary"
	"net"
)

func Handler(nodeID string, f func(query []byte, connection *Connection)) {
	b := sha256.Sum256([]byte(nodeID))
	handler := &handler{
		f: f,
	}
	handler.nodeID[0] = binary.LittleEndian.Uint64(b[:8])
	handler.nodeID[1] = binary.LittleEndian.Uint64(b[8:16])
	handler.nodeID[2] = binary.LittleEndian.Uint64(b[16:24])
	handler.nodeID[3] = binary.LittleEndian.Uint64(b[24:])

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
