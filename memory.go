package node

import (
	"crypto/sha256"
	"crypto/tls"
)

type memory struct {
	tlsConfig *tls.Config
	next      *node
	in        chan *incomingPackage
	drop      chan *struct{}
}

func (m *memory) process() {
	for {
		select {
		case p := <-m.in:
			switch {
			case p.b[0]>>7&1 == 0:
				p.nextMac = sha256.Sum256(p.b[1:p.n])
				if m.next != nil {
					m.next.in <- p
				} else {
					m.next = &node{
						conn: tls.Server(&dataport{}, m.tlsConfig),
						heap: &heap{
							cap: 512,
						},
						in: make(chan *incomingPackage),
					}
					m.next.incoming = p
					m.next.lastIncoming = p
					m.next.initalMac = p.nextMac
					go m.next.process()
				}
			case p.b[0]>>7&1 == 1:
				p.nextMac = sha256.Sum256(p.b[65:p.n])
				if m.next != nil {
					m.next.in <- p
				}
			}
		case <-m.drop:
		}
	}
}
