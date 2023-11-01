package node

import (
	"crypto/sha256"
	"net"
)

func Handler(nodeID string, f func(query []byte, connection *Connection)) {
	handler := &handler{
		nodeID: sha256.Sum256([]byte(nodeID)),
		f:      f,
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
	nodeID [32]byte
	f      func(query []byte, connection *Connection)
}

var (
	initHandlers []*handler
	handlers     []*handler
)
