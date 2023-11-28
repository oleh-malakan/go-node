package node

type heapItem struct {
	incoming  *incomingPackage
	indexNext *heapItem
	indexPrev *heapItem
	next      *heapItem
	prev      *heapItem
}

type heap struct {
	heap *heapItem
	last *heapItem
	len  int
	cap  int // required > 0
}

func (h *heap) put(incoming *incomingPackage) {

}

func (h *heap) find(nextMac []byte) (next *incomingPackage) {

	return nil
}
