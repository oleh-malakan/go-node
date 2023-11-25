package node

import (
	"crypto/sha256"
	"crypto/tls"
)

type controller struct {
	tlsConfig    *tls.Config
	next         *container
	in           chan *incomingPackage
	nextDrop     chan *container
	blackList    *blackList
	limitPrepare chan *struct{}
}

func (c *controller) process() {
	go c.blackList.process()

	go func() {
		for p := range c.in {
			switch {
			case p.b[0]>>7&1 == 0:
				select {
				case c.limitPrepare <- nil:
					p.nextMac = sha256.Sum256(p.b[1:p.n])
					new := &container{
						conn: tls.Server(&dataport{}, c.tlsConfig),
						heap: &heap{
							cap: 512,
						},
						in:        make(chan *incomingPackage),
						nextDrop:  make(chan *container),
						reset:     make(chan *struct{}),
						controler: c,
					}
					new.incoming = p
					new.lastIncoming = p
					go new.prepare()
				default:
				}
			case p.b[0]>>7&1 == 1:
				p.nextMac = sha256.Sum256(p.b[65:p.n])
				c.blackList.in <- p
			}
		}
	}()

	for {
		select {
		case new := <-c.blackList.outPrepare:
			new.next = c.next
			c.next = new
			go new.process()
		case p := <-c.blackList.out:
			if c.next != nil {
				c.next.in <- p
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
