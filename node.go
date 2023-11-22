package node

import (
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"net"
)

type handler struct {
	nodeID [32]byte
	f      func(connection *Connection)
}

var (
	handlers []*handler
)

func Handler(nodeID string, f func(connection *Connection)) {
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

	core := &core{}

	return core.do(handlers, tlsConfig, address, nodeAddresses...)
}
