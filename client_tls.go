package node

import (
	"crypto/tls"
	"errors"
	"net"
)

func DialTLS(tlsConfig *tls.Config, nodeAddresses ...*net.UDPAddr) (*Client, error) {
	if len(nodeAddresses) == 0 {
		return nil, errors.New("node address not specified")
	}

	return nil, nil
}
