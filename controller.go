package node

import (
	"crypto/sha256"
	"crypto/tls"
)

type controller struct {
	tlsConfig *tls.Config
	next      *container
	in        chan *incomingPackage
	nextDrop  chan *container
}

func (c *controller) process() {
	for {
		select {
		case p := <-c.in:
			switch {
			case p.b[0]>>7&1 == 0:
				p.nextMac = sha256.Sum256(p.b[1:p.n])
				new := &container{
					conn: tls.Server(&dataport{}, c.tlsConfig),
					heap: &heap{
						cap: 512,
					},
					in:       make(chan *incomingPackage),
					nextDrop: make(chan *container),
					reset:    make(chan *struct{}),
				}
				new.incoming = p
				new.lastIncoming = p
				new.next = c.next
				c.next = new
				go new.process()
			case p.b[0]>>7&1 == 1:
				p.nextMac = sha256.Sum256(p.b[65:p.n])
				if c.next != nil {
					c.next.in <- p
				}
			}
		case dropNode := <-c.nextDrop:
			c.next = dropNode.next
			if c.next != nil {
				c.next.drop = c.nextDrop
				select {
				case c.next.reset <- nil:
				default:
				}
			}
		}
	}
}
