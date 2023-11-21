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

func TestHeapCap0(t *testing.T) {
	heap := tHeap{}

	heap.Put(arr[0])
	heap.Put(arr[1])
	
	_, last := heap.Find(arr[0].nextMac[0:32])
	if last == nil {
		t.Fatal("next not fount")
	}

	if heap.heap != nil {
		t.Fatal("heap not nil")
	}
}

func TestHeapCap1(t *testing.T) {
	heap := tHeap{
		cap: 1,
	}

	heap.Put(arr[0])
	heap.Put(arr[1])
	
	_, last := heap.Find(arr[0].nextMac[0:32])
	if last == nil {
		t.Fatal("next not fount")
	}

	if heap.heap != nil {
		t.Fatal("heap not nil")
	}
}


func TestHeapCap5(t *testing.T) {
	heap := tHeap{
		cap: 5,
	}

	heap.Put(arr[0])
	heap.Put(arr[3])
	heap.Put(arr[2])
	heap.Put(arr[1])
	heap.Put(arr[4])
	heap.Put(arr[5])
	
	_, last := heap.Find(arr[0].nextMac[0:32])
	if last == nil {
		t.Fatal("next not fount")
	}

	if heap.heap != nil {
		t.Fatal("heap not nil")
	}
}
