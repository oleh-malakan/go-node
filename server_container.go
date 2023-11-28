package node

type serverContainer struct {
	core     *core
	next     *serverContainer
	in       chan *incomingPackage
	nextDrop chan *serverContainer
	drop     chan *serverContainer
	reset    chan *struct{}
	isDrop   bool
}

func (c *serverContainer) process() {
	for !c.isDrop {
		select {
		case i := <-c.in:
			o := c.core.outgoing
			for o != nil {
				if compare16(o.b[1:17], i.b[1:17]) {
					c.core.in(i)
					continue
				}
				o = o.prev
			}

			if c.next != nil {
				c.next.in <- i
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
