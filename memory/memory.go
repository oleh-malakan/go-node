package memory

const (
	cap   = 1000
	depth = 1000000
)

type indexArray[T any] struct {
	array     []*T
	indexFree []int16
	lenFree   int16
}

func (a *indexArray[T]) put(v *T) int16 {
	if v != nil {
		if a.lenFree > 0 {
			a.lenFree--
			index := a.indexFree[a.lenFree]
			a.array[index] = v
			return index
		} else if index := int16(len(a.array)); index < cap {
			a.array = append(a.array, v)
			a.indexFree = append(a.indexFree, 0)
			return index
		}
	}

	return -1
}

func (a *indexArray[T]) get(index int16) *T {
	if int(index) < len(a.array) {
		return a.array[index]
	}

	return nil
}

func (a *indexArray[T]) free(index int16) {
	if lenArray := len(a.array); int(index) < lenArray && a.array[index] != nil {
		if lenArray--; int(index) < lenArray {
			a.array[index] = nil
			a.indexFree[a.lenFree] = index
			a.lenFree++
		} else {
			a.array = a.array[:lenArray]
			a.indexFree = a.indexFree[:lenArray]
		}
		if int(a.lenFree) == len(a.array) {
			a.array = nil
			a.indexFree = nil
			a.lenFree = 0
		}
	}
}

type orderContainer[T any] struct {
	array []*T
}

type page[T any] struct {
	p *orderContainer[*indexArray[T]]
}

func (p *page[T]) open() (int16, *indexArray[T]) {
	return 0, nil
}

func (p *page[T]) get(index int16) *indexArray[T] {
	return nil
}

func (p *page[T]) free(index int16) {

}

type bank[T any] struct {
	b *orderContainer[*page[T]]
}

func (b *bank[T]) open() (int16, *page[T]) {
	return 0, nil
}

func (b *bank[T]) get(index int16) *page[T] {
	return nil
}

func (b *bank[T]) free(index int16) {

}

type Memory[T any] struct {
	bank *bank[T]
}

func (m *Memory[T]) Put(v *T) int32 {
	if z, page := m.bank.open(); page != nil {
		if y, array := page.open(); array != nil {
			return int32(int(z)*depth + int(y)*cap + int(array.put(v)))
		}
	}

	return -1
}

func (m *Memory[T]) Get(index int32) *T {
	xy := int(index % depth)
	x := int16(xy % cap)
	y := int16(xy / cap)
	z := int16(index / depth)

	if page := m.bank.get(z); page != nil {
		if array := page.get(y); array != nil {
			return array.get(x)
		}
	}

	return nil
}

func (m *Memory[T]) Free(index int32) {
}
