package node

import "crypto/sha256"

type memory struct {
	next *node
	in   chan *incomingPackage
	drop chan *struct{}
}

func (m *memory) process() {
	for {
		select {
		case p := <-m.in:
			if m.next != nil {
				switch {
				case p.b[0]>>7&1 == 0:
					p.nextMac = sha256.Sum256(p.b[1:p.n])
				case p.b[0]>>7&1 == 1:
					p.nextMac = sha256.Sum256(p.b[65:p.n])
				}

				m.next.in <- p
			} else {

			}
		case <-m.drop:
		}
	}
}
