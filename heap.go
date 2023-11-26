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
	item := &heapItem{
		incoming: incoming,
	}

LOOP:
	if h.last != nil {
		if h.cap <= h.len {
			h.len--

			if h.last.indexNext != nil {
				h.last.indexNext.indexPrev = nil
			}
			if h.last.indexPrev != nil {
				h.last.indexPrev.indexNext = nil
			}

			if h.last.prev != nil {
				h.last = h.last.prev
				h.last.next = nil
			} else {
				h.last = nil
				h.heap = nil
				goto LOOP
			}
		}

		heap := h.heap
		for heap != nil && (item.indexPrev == nil || item.indexNext == nil) {
			if compareID(incoming.b[33:65], heap.incoming.b[33:65]) {
				return
			}
			if item.indexNext == nil && compareID(incoming.nextMac[0:32], heap.incoming.b[33:65]) {
				item.indexNext = heap
			}
			if item.indexPrev == nil && compareID(heap.incoming.nextMac[0:32], incoming.b[33:65]) {
				item.indexPrev = heap
				if item.indexPrev.indexNext != nil {
					return
				}
			}

			heap = heap.next
		}

		if heap != nil {
			if heap.next != nil {
				heap.next.prev = item
				item.next = heap.next
			}
			heap.next = item
			item.prev = heap
			if heap.next == nil {
				h.last = heap
			}
		} else {
			item.prev = h.last
			h.last.next = item
			h.last = item
		}

		if item.indexPrev != nil {
			item.indexPrev.indexNext = item
		}
	} else {
		h.heap = item
		h.last = item
	}

	h.len++
}

func (h *heap) find(nextMac []byte) (next, last *incomingPackage) {
	delete := func(item *heapItem) {
		if item != nil {
			if item.prev != nil {
				item.prev.next = item.next
				if item.next != nil {
					item.next.prev = item.prev
				} else {
					h.last = item.prev
				}
			} else {
				h.heap = item.next
				if h.heap != nil {
					h.heap.prev = nil
				} else {
					h.last = nil
				}
			}
			h.len--
		}
	}

	heap := h.heap
	for heap != nil {
		if compareID(nextMac, heap.incoming.b[33:65]) {
			next = heap.incoming
			last = next
			delete(heap)
			index := heap.indexNext
			for index != nil {
				last.next = index.incoming
				last = last.next
				delete(index)
				index = index.indexNext
			}
		}
		heap = heap.next
	}

	return
}
