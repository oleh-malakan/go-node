package node

import "crypto/tls"

type node struct {
	tlsConfig           *tls.Config
	conn                *tls.Conn
	lastIncomingPackage *incomingPackage
	incomingPackage     *incomingPackage
	outgoingPackage     *outgoingPackage
	initalMac           [32]byte
	heap                *heap
	next                *node
	in                  chan *incomingPackage
	drop                chan *struct{}
	isDrop              bool
}

func (n *node) process() {
	for {
		select {
		case p := <-n.in:
			if !n.isDrop {
				switch {
				case p.b[0]>>7&1 == 0:
					if compareID(n.initalMac[0:32], p.nextMac[0:32]) {
						continue
					} else if n.next == nil {
						n.next = &node{
							conn: tls.Server(&dataport{}, n.tlsConfig),
							heap: &heap{
								cap: 512,
							},
						}
						n.next.incomingPackage = p
						n.next.lastIncomingPackage = p
						n.next.initalMac = p.nextMac
						go n.next.process()

						continue
					} else {
						n.next.in <- p
					}

				case p.b[0]>>7&1 == 1:
					w := n.outgoingPackage
					for w != nil {
						if compareID(w.mac[0:32], p.b[1:33]) {
							if compareID(n.lastIncomingPackage.nextMac[0:32], p.b[33:65]) {
								n.lastIncomingPackage.next = p
								n.lastIncomingPackage = p
								var last *incomingPackage
								n.lastIncomingPackage.next, last = n.heap.find(p.nextMac[0:32])
								if last != nil {
									n.lastIncomingPackage = last
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
			}
		case <-n.drop:
		}
	}
}
