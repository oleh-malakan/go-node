package node

import "crypto/tls"

type container struct {
	conn         *tls.Conn
	lastIncoming *incomingPackage
	incoming     *incomingPackage
	outgoing     *outgoingPackage
	heap         *heap
	next         *container
	in           chan *incomingPackage
	nextDrop     chan *container
	drop         chan *container
	reset        chan *struct{}
	isDrop       bool
}

func (c *container) process() {
	for !c.isDrop {
		select {
		case p := <-c.in:
			w := c.outgoing
			for w != nil {
				if compareID(w.mac[0:32], p.b[1:33]) {
					if compareID(c.lastIncoming.nextMac[0:32], p.b[33:65]) {
						c.lastIncoming.next = p
						c.lastIncoming = p
						var last *incomingPackage
						c.lastIncoming.next, last = c.heap.find(p.nextMac[0:32])
						if last != nil {
							c.lastIncoming = last
						}
					} else {
						c.heap.put(p)
					}

					continue
				}
				w = w.prev
			}

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

	for {
		select {
		case p := <-c.in:
			if c.next != nil {
				c.next.in <- p
			}
		case <-c.reset:
		case c.drop <- c:
			return
		}
	}
}
