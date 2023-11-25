package node

type blackList struct {
	in         chan *incomingPackage
	out        chan *incomingPackage
	inPrepare  chan *container
	outPrepare chan *container
	black      chan [32]byte
}

func (b *blackList) process() {
	for {
		select {
		case p := <-b.in:
			b.out <- p
		case c := <-b.inPrepare:
			b.outPrepare <- c
		case <-b.black:
		}
	}
}
