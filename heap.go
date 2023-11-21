package node

import "time"

type tHeapReadData struct {
	readData *tReadData
	prev     *tHeapItem
	next     *tHeapItem
	time     int64
	timeout  int64
}

type tHeapItem struct {
	data *tHeapReadData
	next *tHeapItem
}

type tHeap struct {
	heap    *tHeapItem
	len     int
	cap     int
	timeout int64
}

func (t *tHeap) Put(r *tReadData) {
	var (
		prev *tHeapItem
		next *tHeapItem
	)
	heap := t.heap
	for heap != nil && (prev == nil || next == nil) {
		if compareID(r.b[33:65], heap.data.readData.b[33:65]) {
			heap.data.time = time.Now().UnixNano()
			return
		}
		if heap.data.prev == nil && compareID(r.nextMac[0:32], heap.data.readData.b[33:65]) {
			next = heap
		}
		if heap.data.next == nil && compareID(heap.data.readData.nextMac[0:32], r.b[33:65]) {
			prev = heap
		}
		heap = heap.next
	}

	heapItem := &tHeapItem{
		data: &tHeapReadData{
			readData: r,
			prev:     prev,
			next:     next,
			time:     time.Now().UnixNano(),
			timeout:  t.timeout,
		},
	}

	if heap.next != nil {

	} else {
		heap.next = heapItem
	}

}

func (t *tHeap) Find(nextMac []byte) (next, last *tReadData) {
	for i := 0; i < len(t.index) && next == nil; i++ {
		if compareID(nextMac, t.heap[t.index[i]].readData.b[33:65]) {
			index := t.index[i]
			next = t.heap[index].readData
			last = next
			t.heap[index] = nil
			t.indexPrepareFree(index)
			for index = t.heap[index].next; index >= 0; {
				last.next = t.heap[index].readData
				last = last.next
				t.heap[t.index[i]] = nil
				t.indexPrepareFree(index)
			}
			t.indexFree()
		}
	}

	return
}
