package node

import (
	"crypto/sha256"
	"crypto/tls"
)

type controller struct {
	config    *Config
	tlsConfig *tls.Config
	next      *container
	in        chan *incomingPackage
	nextDrop  chan *container
}

func (c *controller) process() {
	for {
		select {
		case i := <-c.in:
			switch {
			case i.b[0]>>7&1 == 0:
				new := &container{
					core: &core{
						iPKey: sha256.Sum256(i.b[1:i.n]),
						heap: &heap{
							cap: c.config.HeapCap,
						},
					},
					in:       make(chan *incomingPackage),
					nextDrop: make(chan *container),
					reset:    make(chan *struct{}),
				}
				new.core.conn = tls.Server(new.core, c.tlsConfig)
				new.core.incoming = i
				new.core.lastIncoming = i
				new.next = c.next
				c.next = new
				go new.process()
			case i.b[0]>>7&1 == 1:
				if c.next != nil {
					c.next.in <- i
				}
			}
		case d := <-c.nextDrop:
			c.next = d.next
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
