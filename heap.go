package node

import "time"

type tHeapItem struct {
	readData *tReadData
	prev     int
	next     int
	time     int64
	timeout  int64
}

type tHeap struct {
	heap             []*tHeapItem
	index            []int
	freeIndex        []int
	prepareFreeIndex []int
	cap              int
	timeout          int64
}

func (t *tHeap) Put(r *tReadData) {
	prev := -1
	next := -1
	for i := 0; i < len(t.index) && (prev < 0 || next < 0); i++ {
		if compareID(r.b[33:65], t.heap[t.index[i]].readData.b[33:65]) {
			t.heap[t.index[i]].time = time.Now().UnixNano()
			return
		}
		if t.heap[t.index[i]].prev < 0 {
			if compareID(r.nextMac[0:32], t.heap[t.index[i]].readData.b[33:65]) {
				next = t.index[i]
			}
		}
		if t.heap[t.index[i]].next < 0 {
			if compareID(t.heap[t.index[i]].readData.nextMac[0:32], r.b[33:65]) {
				prev = t.index[i]
			}
		}
	}

	t.heap[t.indexNext()] = &tHeapItem{
		readData: r,
		prev:     prev,
		next:     next,
		time:     time.Now().UnixNano(),
		timeout:  t.timeout,
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

func (t *tHeap) indexNext() (index int) {
	if len(t.freeIndex) > 0 {
		index = t.freeIndex[0]
		t.freeIndex = t.freeIndex[1:len(t.freeIndex)]
		t.indexPut(index)
	} else if len(t.heap) == 0 {
		t.heap = append(t.heap, nil)
		t.indexPut(index)
	} else if index = len(t.heap); index < t.cap {
		t.heap = append(t.heap, nil)
		t.indexPut(index)
	} else {
		index = len(t.heap) - 1
		if t.heap[index].next >= 0 {
			t.heap[t.heap[index].next].prev = -1
		}
		if t.heap[index].prev >= 0 {
			t.heap[t.heap[index].prev].next = -1
		}
	}

	return
}

func (t *tHeap) indexPut(index int) {
	t.index = t.arrayPutIndex(t.index, index)
}

func (t *tHeap) indexPrepareFree(index int) {
	t.prepareFreeIndex = t.arrayPutIndex(t.prepareFreeIndex, index)
}

func (t *tHeap) indexFree() {

	t.prepareFreeIndex = nil
}

func (t *tHeap) arrayPutIndex(arr []int, index int) []int {
	var min int
	for min = 0; min < len(arr) && arr[min] < index; min++ {
	}

	var new []int
	new = append(new, arr[:min]...)
	new = append(new, index)
	new = append(new, arr[min:]...)

	return new
}
