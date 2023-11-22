package node

import (
	"crypto/sha256"
	"encoding/binary"
	"testing"
)

const (
	capHeap = 512
)

var (
	arr []*tReadData
)

func init() {
	var (
		nextMac [32]byte
	)
	for i := 0; i < capHeap; i++ {
		readData := &tReadData{
			b: make([]byte, 1432),
		}
		data := make([]byte, 8)
		binary.BigEndian.PutUint64(data, uint64(i))
		copy(readData.b[65:73], data)
		copy(readData.b[33:65], nextMac[0:32])
		nextMac = sha256.Sum256(readData.b[65:])
		readData.nextMac = nextMac
		arr = append(arr, readData)
	}
}

func heapCheckResult(next *tReadData, offset int, len int) bool {
	for i := offset; i < offset+len; i++ {
		if next != nil {
			v := binary.BigEndian.Uint64(next.b[65:73])
			if int(v) != i {
				return false
			}
		} else {
			return false
		}
		next = next.next
	}

	return true
}

func TestHeapCap0(t *testing.T) {
	heap := tHeap{}

	heap.put(arr[0])
	heap.put(arr[1])

	next, last := heap.find(arr[0].nextMac[0:32])
	if last == nil {
		t.Fatal("next not fount")
	}

	if !heapCheckResult(next, 1, 1) {
		t.Fatal("next failed")
	}

	if heap.heap != nil {
		t.Fatal("heap not nil")
	}
}

func TestHeapCap1(t *testing.T) {
	heap := tHeap{
		cap: 1,
	}

	heap.put(arr[0])
	heap.put(arr[1])

	next, last := heap.find(arr[0].nextMac[0:32])
	if last == nil {
		t.Fatal("next not fount")
	}

	if !heapCheckResult(next, 1, 1) {
		t.Fatal("next failed")
	}

	if heap.heap != nil {
		t.Fatal("heap not nil")
	}
}

func TestHeapCap5(t *testing.T) {
	heap := tHeap{
		cap: 5,
	}

	heap.put(arr[1])
	heap.put(arr[4])
	heap.put(arr[3])
	heap.put(arr[2])
	heap.put(arr[5])
	heap.put(arr[6])

	next, last := heap.find(arr[0].nextMac[0:32])
	if last == nil {
		t.Fatal("next not fount")
	}

	if !heapCheckResult(next, 1, 4) {
		t.Fatal("next failed")
	}

	heap.put(arr[5])

	next, last = heap.find(arr[4].nextMac[0:32])
	if last == nil {
		t.Fatal("next not fount")
	}

	if !heapCheckResult(next, 5, 2) {
		t.Fatal("next failed")
	}

	if heap.heap != nil {
		t.Fatal("heap not nil")
	}
}
