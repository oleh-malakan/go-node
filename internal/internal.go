package internal

type IndexArray[T any] struct {
	array []*T
}

func (i *IndexArray[T]) Get(index int) *T {
	if index >= 0 && index < len(i.array) {
		return i.array[index]
	}

	return nil
}

func (i *IndexArray[T]) Put(v *T) (index int) {
	return
}

func (i *IndexArray[T]) Free(index int) {

}

type CIDManager struct{}

func (c *CIDManager) CID() []byte {
	var ID []byte

	return ID[:]
}

func (c *CIDManager) Put(cid []byte) bool {
	return true
}

func (i *IndexArray[T]) Free(index int) {

}
