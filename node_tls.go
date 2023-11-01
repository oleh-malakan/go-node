package node

import (
	"crypto/tls"
	"net"
)

func DoTLS(tlsConfig *tls.Config, address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) error {

	return nil
}
