package internal

type IndexArray[T any] struct {
	array []*T
}

func (i *IndexArray[T]) Get(index int64) *T {
	platformIndex := int(index)
	if platformIndex >= 0 && platformIndex < len(i.array) {
		return i.array[platformIndex]
	}

	return nil
}

func (i *IndexArray[T]) Put(v *T) (index int64) {
	return
}

func (i *IndexArray[T]) Free(index int64) {

}
