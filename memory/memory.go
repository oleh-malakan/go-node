package memory

const (
	cap   = 1000
	depth = 1000000
)

type orderVector[T any] struct {
	array     []*T
	indexFree []int16
}

func (v *orderVector[T]) put(value *T) int {
	if value != nil {
		if len(v.indexFree) > 0 {
			index := v.indexFree[0]
			v.indexFree = v.indexFree[1:]
			v.array[index] = value
			return int(index)
		} else if index := len(v.array); index < cap {
			v.array = append(v.array, value)
			return index
		}
	}

	return -1
}

func (v *orderVector[T]) get(index int) *T {
	if index < len(v.array) {
		return v.array[index]
	}

	return nil
}

func (v *orderVector[T]) free(index int) {
	if index < len(v.array) && v.array[index] != nil {
		v.array[index] = nil
		for i := 0; i < len(v.indexFree); i++ {
			if index < int(v.indexFree[i]) {
				temp := v.indexFree[i]
				v.indexFree[i] = int16(index)
				index = int(temp)
			}
		}
		v.indexFree = append(v.indexFree, int16(index))

		for len(v.indexFree) > 0 && int(v.indexFree[len(v.indexFree)-1]) == len(v.array)-1 {
			v.array = v.array[:len(v.array)-1]
			v.indexFree = v.indexFree[:len(v.indexFree)-1]
		}
	}
}

func (v *orderVector[T]) len() int {
	return len(v.array)
}

type column[T any] struct {
	array     []*T
	indexFree []int16
	lenFree   int16
}

func (c *column[T]) put(value *T) int {
	if value != nil {
		if c.lenFree > 0 {
			c.lenFree--
			index := c.indexFree[c.lenFree]
			c.array[index] = value
			return int(index)
		} else if index := len(c.array); index < cap {
			c.array = append(c.array, value)
			c.indexFree = append(c.indexFree, 0)
			return index
		}
	}

	return -1
}

func (c *column[T]) get(index int) *T {
	if index < len(c.array) {
		return c.array[index]
	}

	return nil
}

func (c *column[T]) free(index int) {
	if lenArray := len(c.array); index < lenArray && c.array[index] != nil {
		if lenArray--; int(index) < lenArray {
			c.array[index] = nil
			c.indexFree[c.lenFree] = int16(index)
			c.lenFree++
		} else {
			c.array = c.array[:lenArray]
			c.indexFree = c.indexFree[:lenArray]
		}
		if int(c.lenFree) == len(c.array) {
			c.array = nil
			c.indexFree = nil
			c.lenFree = 0
		}
	}
}

func (c *column[T]) len() int {
	return len(c.array)
}

type page[T any] struct {
	vector *orderVector[column[T]]
	cursor int16
}

func (p *page[T]) open() (int, *column[T]) {
	var cursor int
	if p.cursor > 0 {
		cursor = int(p.cursor)
	}

	for i := cursor; i < p.vector.len(); i++ {
		if column := p.vector.get(i); column != nil && column.len() < cap {
			return i, column
		}
	}

	if p.vector.len() < cap {
		column := &column[T]{}
		return p.vector.put(column), column
	}

	return 0, nil
}

func (p *page[T]) get(index int) *column[T] {
	return nil
}

func (p *page[T]) free(index int) {

}

func (p *page[T]) len() int {
	return p.vector.len()
}

type bank[T any] struct {
	vector *orderVector[page[T]]
	cursor int16
}

func (b *bank[T]) open() (int, *page[T]) {
	var cursor int
	if b.cursor > 0 {
		cursor = int(b.cursor)
	}

	for i := cursor; i < b.vector.len(); i++ {
		if page := b.vector.get(i); page != nil && page.len() < cap {
			return i, page
		}
	}

	if b.vector.len() < cap {
		page := &page[T]{}
		return b.vector.put(page), page
	}

	return 0, nil
}

func (b *bank[T]) get(index int) *page[T] {
	return nil
}

func (b *bank[T]) free(index int) {

}

type Memory[T any] struct {
	bank *bank[T]
}

func (m *Memory[T]) Put(v *T) int32 {
	if z, page := m.bank.open(); page != nil {
		if y, column := page.open(); column != nil {
			return int32(int(z)*depth + int(y)*cap + int(column.put(v)))
		}
	}

	return -1
}

func (m *Memory[T]) Get(index int32) *T {
	xy := int(index % depth)
	x := xy % cap
	y := xy / cap
	z := int(index / depth)

	if page := m.bank.get(z); page != nil {
		if column := page.get(y); column != nil {
			return column.get(x)
		}
	}

	return nil
}

func (m *Memory[T]) Free(index int32) {
}
