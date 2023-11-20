package node

import "testing"

func TestArrayPutIndex(t *testing.T) {
	heap := &tHeap{}

	arr := []int{}
	arr = heap.arrayPutIndex(arr, 2)
	res := []int{2}
	if !checkResult(arr, res) {
		t.Fatalf("\n1: \narr: %v\nres: %v", arr, res)
	}

	arr = []int{2, 3, 6, 7, 9}
	arr = heap.arrayPutIndex(arr, 5)
	res = []int{2, 3, 5, 6, 7, 9}
	if !checkResult(arr, res) {
		t.Fatalf("\n2: \narr: %v\nres: %v", arr, res)
	}

	arr = heap.arrayPutIndex(arr, 10)
	res = []int{2, 3, 5, 6, 7, 9, 10}
	if !checkResult(arr, res) {
		t.Fatalf("\n3: \narr: %v\nres: %v", arr, res)
	}

	arr = heap.arrayPutIndex(arr, 0)
	res = []int{0, 2, 3, 5, 6, 7, 9, 10}
	if !checkResult(arr, res) {
		t.Fatalf("\n4: \narr: %v\nres: %v", arr, res)
	}
}

func checkResult(a, b []int) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
