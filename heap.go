package node

import "time"

type tHeapItem struct {
	readData  *tReadData
	nextIndex *tHeapItem
	time      int64
	timeout   int64
	next      *tHeapItem
	prev      *tHeapItem
}

type tHeap struct {
	heap    *tHeapItem
	len     int
	cap     int
	timeout int64
}

func (t *tHeap) Put(r *tReadData) {
	var (
		prevIndex *tHeapItem
		nextIndex *tHeapItem
	)
	heap := t.heap
	for heap != nil && (prevIndex == nil || nextIndex == nil) {
		if compareID(r.b[33:65], heap.readData.b[33:65]) {
			heap.time = time.Now().UnixNano()
			return
		}
		if heap.nextIndex == nil && compareID(r.nextMac[0:32], heap.readData.b[33:65]) {
			nextIndex = heap
		}
		if prevIndex == nil && compareID(heap.readData.nextMac[0:32], r.b[33:65]) {
			prevIndex = heap
			if prevIndex.nextIndex != nil {
				prevIndex.nextIndex.time = time.Now().UnixNano()
				return	
			}
		}

		heap = heap.next
	}

	heapItem := &tHeapItem{
		readData:  r,
		nextIndex: nextIndex,
		time:      time.Now().UnixNano(),
		timeout:   t.timeout,
	}


	if t.cap <= t.len {

	}

	if prevIndex != nil {
		prevIndex.nextIndex = heapItem
	}

	if heap != nil {
		if heap.next != nil {
			heap.next.prev = heapItem
			heapItem.next = heap.next
		}
		heap.next = heapItem
		heapItem.prev = heap
		t.len++
	} else {

	}
}

func (t *tHeap) Find(nextMac []byte) (next, last *tReadData) {
	heap := t.heap
	for heap != nil {
		if compareID(nextMac, heap.readData.b[33:65]) {

			index := heap
			next = index.readData
			last = next

			//t.heap[index] = nil
			for index = t.heap[index].next; index >= 0; {
				last.next = t.heap[index].readData
				last = last.next
				//t.heap[t.index[i]] = nil
			}
		}
	}

	return
}
