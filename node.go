package node

import "crypto/tls"

type node struct {
	tlsConfig    *tls.Config
	conn         *tls.Conn
	lastIncoming *incomingPackage
	incoming     *incomingPackage
	outgoing     *outgoingPackage
	initalMac    [32]byte
	heap         *heap
	next         *node
	in           chan *incomingPackage
	nextDrop     chan *node
	drop         chan *node
	reset        chan *struct{}
	isDrop       bool
}

func (n *node) process() {
	for !n.isDrop {
		select {
		case p := <-n.in:
			switch {
			case p.b[0]>>7&1 == 0:
				if compareID(n.initalMac[0:32], p.nextMac[0:32]) {
					continue
				} else if n.next != nil {
					n.next.in <- p
				} else {
					n.next = newNode(p, n.nextDrop, n.tlsConfig)
					go n.next.process()

					continue
				}
			case p.b[0]>>7&1 == 1:
				w := n.outgoing
				for w != nil {
					if compareID(w.mac[0:32], p.b[1:33]) {
						if compareID(n.lastIncoming.nextMac[0:32], p.b[33:65]) {
							n.lastIncoming.next = p
							n.lastIncoming = p
							var last *incomingPackage
							n.lastIncoming.next, last = n.heap.find(p.nextMac[0:32])
							if last != nil {
								n.lastIncoming = last
							}
						} else {
							n.heap.put(p)
						}

						continue
					}
					w = w.prev
				}

				if n.next != nil {
					n.next.in <- p
				}
			}
		case dropNode := <-n.nextDrop:
			n.next = dropNode.next
			if n.next != nil {
				n.next.drop = n.nextDrop
				select {
				case n.next.reset <- nil:
				default:
				}
			}
		}
	}

	for {
		select {
		case p := <-n.in:
			if n.next != nil {
				n.next.in <- p
			}
		case <-n.reset:
		case n.drop <- n:
			return
		}
	}
}

func newNode(incoming *incomingPackage, nextDrop chan *node, tlsConfig *tls.Config) *node {
	new := &node{
		conn: tls.Server(&dataport{}, tlsConfig),
		heap: &heap{
			cap: 512,
		},
		in:       make(chan *incomingPackage),
		nextDrop: make(chan *node),
		reset:    make(chan *struct{}),
	}
	new.incoming = incoming
	new.lastIncoming = incoming
	new.initalMac = incoming.nextMac
	new.drop = nextDrop

	return new
}
