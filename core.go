package node

import "crypto/tls"

type core struct {
	conn         *tls.Conn
	lastIncoming *incomingPackage
	incoming     *incomingPackage
	outgoing     *outgoingPackage
	heap         *heap
}

func (c *core) in(incoming *incomingPackage) {
	if compareID(c.lastIncoming.nextMac[0:32], incoming.b[33:65]) {
		c.lastIncoming.next = incoming
		c.lastIncoming = incoming
		var last *incomingPackage
		c.lastIncoming.next, last = c.heap.find(incoming.nextMac[0:32])
		if last != nil {
			c.lastIncoming = last
		}
	} else {
		c.heap.put(incoming)
	}
}

func (c *core) process() {
	select {
	default:
	}
}
