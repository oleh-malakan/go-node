package node

import "time"

type container struct {
	core     *core
	next     *container
	in       chan *incomingPackage
	nextDrop chan *container
	drop     chan *container
	reset    chan *struct{}
	isDrop   bool
}

func (c *container) do() {
	for !c.isDrop {
		select {
		case i := <-c.in:
			o := c.core.outgoing
			for o != nil {
				if compareID(o.mac[0:32], i.b[1:33]) {
					c.core.in(i)
					continue
				}
				o = o.prev
			}

			if c.next != nil {
				c.next.in <- i
			}
		default:
			c.core.process()
			time.Sleep(0)
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

	for {
		select {
		case i := <-c.in:
			if c.next != nil {
				c.next.in <- i
			}
		case <-c.reset:
		case c.drop <- c:
			return
		}
	}
}
