package node

import (
	"encoding/binary"
	"testing"
)

func heapCheckResult(next *incomingPackage, offset int, len int) bool {
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
	heap := heap{}

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
	heap := heap{
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
	heap := heap{
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

func TestHeap1(t *testing.T) {
	heap := heap{
		cap: 512,
	}

	if heap.heap != nil {
		t.Fatal("heap not nil")
	}
}

func TestHeap2(t *testing.T) {
	heap := heap{
		cap: 512,
	}

	if heap.heap != nil {
		t.Fatal("heap not nil")
	}
}
