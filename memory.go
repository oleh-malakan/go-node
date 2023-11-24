package node

import (
	"crypto/sha256"
	"crypto/tls"
)

type memory struct {
	tlsConfig *tls.Config
	next      *node
	in        chan *incomingPackage
	nextDrop  chan *node
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
					m.next = newNode(p, m.nextDrop, m.tlsConfig)
					go m.next.process()
				}
			case p.b[0]>>7&1 == 1:
				p.nextMac = sha256.Sum256(p.b[65:p.n])
				if m.next != nil {
					m.next.in <- p
				}
			}
		case dropNode := <-m.nextDrop:
			m.next = dropNode.next
			if m.next != nil {
				m.next.drop = m.nextDrop
				select {
				case m.next.reset <- nil:
				default:
				}
			}
		}
	}
}
