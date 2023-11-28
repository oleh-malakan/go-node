package node

type heapItem struct {
	incoming *incomingPackage
	next     *heapItem
	prev     *heapItem
}

type heap struct {
	heap *heapItem
	last *heapItem
	len  int
	cap  int
}

func (h *heap) put(incoming *incomingPackage) {
	cur := h.heap
	for cur != nil {
		if compare16(cur.incoming.b[17:33], incoming.b[17:33]) {
			return
		}
	}

	if h.cap <= h.len {
		if h.heap != nil {
			h.heap = h.heap.next
			if h.heap == nil {
				h.last = nil
			}
			h.len--
		}
	}

	item := &heapItem{
		incoming: incoming,
	}
	if h.last != nil {
		item.prev = h.last
		h.last.next = item
		h.last = item
	} else {
		h.heap = item
		h.last = item
	}
	h.len++
}

func (h *heap) find(next []byte) *incomingPackage {
	cur := h.heap
	for cur != nil {
		if compare8(next, cur.incoming.b[17:25]) {
			if cur.prev != nil {
				cur.prev.next = cur.next
			} else {
				h.heap = cur.next
				if h.heap == nil {
					h.last = nil
				}
			}

			return cur.incoming
		}

		cur = cur.next
	}

	return nil
}
