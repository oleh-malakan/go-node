package memory

const (
	cap   = 1000
	depth = 1000000
)

type indexVector[T any] struct {
	array     []*T
	indexFree []int16
	lenFree   int16
}

func (v *indexVector[T]) put(value *T) int16 {
	if value != nil {
		if v.lenFree > 0 {
			v.lenFree--
			index := v.indexFree[v.lenFree]
			v.array[index] = value
			return index
		} else if index := int16(len(v.array)); index < cap {
			v.array = append(v.array, value)
			v.indexFree = append(v.indexFree, 0)
			return index
		}
	}

	return -1
}

func (v *indexVector[T]) get(index int16) *T {
	if int(index) < len(v.array) {
		return v.array[index]
	}

	return nil
}

func (v *indexVector[T]) free(index int16) {
	if lenArray := len(v.array); int(index) < lenArray && v.array[index] != nil {
		if lenArray--; int(index) < lenArray {
			v.array[index] = nil
			v.indexFree[v.lenFree] = index
			v.lenFree++
		} else {
			v.array = v.array[:lenArray]
			v.indexFree = v.indexFree[:lenArray]
		}
		if int(v.lenFree) == len(v.array) {
			v.array = nil
			v.indexFree = nil
			v.lenFree = 0
		}
	}
}

func (v *indexVector[T]) len() int16 {
	return int16(len(v.array))
}

type orderVector[T any] struct {
	array     []*T
	indexFree []int16
}

func (v *orderVector[T]) put(value *T) int16 {
	if value != nil {
		if len(v.indexFree) > 0 {

		} else {

		}
	}

	return -1
}

func (v *orderVector[T]) get(index int16) *T {
	if int(index) < len(v.array) {
		return v.array[index]
	}

	return nil
}

func (v *orderVector[T]) free(index int16) {
}

func (v *orderVector[T]) len() int16 {
	return int16(len(v.array))
}

type page[T any] struct {
	p *orderVector[*indexVector[T]]
}

func (p *page[T]) open() (int16, *indexVector[T]) {
	return 0, nil
}

func (p *page[T]) get(index int16) *indexVector[T] {
	return nil
}

func (p *page[T]) free(index int16) {

}

type bank[T any] struct {
	b *orderVector[*page[T]]
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
