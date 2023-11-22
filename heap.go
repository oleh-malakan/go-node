package node

import "time"

type tHeapItem struct {
	readData  *tReadData
	indexNext *tHeapItem
	indexPrev *tHeapItem
	time      int64
	timeout   int64
	next      *tHeapItem
	prev      *tHeapItem
}

type tHeap struct {
	heap    *tHeapItem
	last    *tHeapItem
	len     int
	cap     int // required > 0
	timeout int64
}

func (t *tHeap) put(r *tReadData) {
	if t.cap <= t.len {
		t.len--
		if t.last.prev != nil {
			t.last = t.last.prev
			t.last.next = nil
		} else {
			t.last = nil
		}
	}
	
	var (
		indexPrev *tHeapItem
		indexNext *tHeapItem
	)
	heap := t.heap
	for heap != nil && (indexPrev == nil || indexNext == nil) {
		if compareID(r.b[33:65], heap.readData.b[33:65]) {
			heap.time = time.Now().UnixNano()
			return
		}
		if indexNext == nil && compareID(r.nextMac[0:32], heap.readData.b[33:65]) {
			indexNext = heap
		}
		if indexPrev == nil && compareID(heap.readData.nextMac[0:32], r.b[33:65]) {
			indexPrev = heap
			if indexPrev.indexNext != nil {
				indexPrev.indexNext.time = time.Now().UnixNano()
				return
			}
		}

		heap = heap.next
	}

	heapItem := &tHeapItem{
		readData:  r,
		indexNext: indexNext,
		indexPrev: indexPrev,
		time:      time.Now().UnixNano(),
		timeout:   t.timeout,
	}

	if t.last != nil {
		if heap != nil {
			if heap.next != nil {
				heap.next.prev = heapItem
				heapItem.next = heap.next
			}
			heap.next = heapItem
			heapItem.prev = heap
			if heap.next == nil {
				t.last = heap
			}
		} else {
			heapItem.prev = t.last
			t.last.next = heapItem
			t.last = heapItem
		}

		if indexPrev != nil {
			indexPrev.indexNext = heapItem
		}
	} else {
		t.heap = heapItem
		t.last = heapItem
	}

	t.len++
}

func (t *tHeap) find(nextMac []byte) (next, last *tReadData) {
	delete := func(heapItem *tHeapItem) {
		if heapItem != nil {
			if heapItem.prev != nil {
				heapItem.prev.next = heapItem.next
				if heapItem.next == nil {
					t.last = heapItem.prev
				}
			} else {
				t.heap = heapItem.next
				if t.heap == nil {
					t.last = nil
				}
			}
			t.len--
		}
	}

	heap := t.heap
	for heap != nil {
		if compareID(nextMac, heap.readData.b[33:65]) {
			next = heap.readData
			last = next
			delete(heap)
			index := heap.indexNext
			for index != nil {
				last.next = index.readData
				last = last.next
				delete(index)
				index = index.indexNext
			}
		}
		heap = heap.next
	}

	return
}
