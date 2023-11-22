package node

import (
	"crypto/sha256"
	"crypto/tls"
	"net"
	"time"
)

type client struct {
	conn         *tls.Conn
	lastReadData *readData
	readData     *readData
	heap         *heap
	writeData    *writeData
	next         *client
	lock         chan *struct{}
	drop         bool
}

type core struct {
	memory     *client
	memoryLock chan *struct{}
	tlsConfig  *tls.Config
}

func (c *core) do(handlers []*handler, tlsConfig *tls.Config, address *net.UDPAddr, nodeAddresses ...*net.UDPAddr) error {
	c.tlsConfig = tlsConfig

	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		return err
	}

	c.memoryLock = make(chan *struct{}, 1)
	for {
		r := &readData{
			b: make([]byte, 1432),
		}
		r.n, r.rAddr, r.err = conn.ReadFromUDP(r.b)

		go c.bypass(r)
	}
}

func (c *core) bypass(r *readData) {
	switch {
	case r.b[0]>>7&1 == 0:
		r.nextMac = sha256.Sum256(r.b[1:r.n])
		client := &client{
			conn: tls.Server(&dataport{}, c.tlsConfig),
			lock: make(chan *struct{}, 1),
			heap: &heap{
				cap:     512,
				timeout: int64(time.Duration(50) * time.Millisecond),
			},
		}
		client.readData = r
		client.lastReadData = r
		client.drop = false

		c.memoryLock <- nil
		if c.memory != nil {
			client.next = c.memory
			c.memory = client
		} else {
			c.memory = client
		}
		<-c.memoryLock
	case r.b[0]>>7&1 == 1 && c.memory != nil:
		r.nextMac = sha256.Sum256(r.b[65:r.n])
		var current *client
		c.memoryLock <- nil
		if c.memory != nil {
			current = c.memory
			current.lock <- nil
		}
		<-c.memoryLock
		for current != nil {
			var next *client
			if !current.drop {
				w := current.writeData
				for w != nil {
					if compareID(w.mac[0:32], r.b[1:33]) {
						if compareID(current.lastReadData.nextMac[0:32], r.b[33:65]) {
							current.lastReadData.next = r
							current.lastReadData = r
							var last *readData
							current.lastReadData.next, last = current.heap.find(r.nextMac[0:32])
							if last != nil {
								current.lastReadData = last
							}
						} else {
							current.heap.put(r)
						}

						next = nil
						goto FOUND
					}
					w = w.prev
				}
			}

			next = current.next
			if next != nil {
				next.lock <- nil
			}
		FOUND:
			<-current.lock
			current = next
		}
	}
}
