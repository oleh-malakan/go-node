package node

import (
	"crypto/sha256"
	"encoding/binary"
)

const (
	capHeap = 512
)

var (
	arr []*incomingPackage
)

func init() {
	var (
		nextMac [32]byte
	)
	for i := 0; i < capHeap; i++ {
		readData := &incomingPackage{
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
