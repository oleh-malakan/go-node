package memory

const (
	capRow    = 1000000
	capColumn = 1000
)

type column[T any] struct {
	row       []*T
	indexFree []int32
	lenFree   int
}

type Memory[T any] struct {
	column    []*column[T]
	indexFree []int16
	cursor    int
}

func (m *Memory[T]) Put(v *T) int32 {
	for m.cursor < len(m.column) {
		if m.column[m.cursor] != nil && len(m.column[m.cursor].row) < capColumn {
			if m.column[m.cursor].lenFree > 0 {
				m.column[m.cursor].lenFree--
				i := m.column[m.cursor].indexFree[m.column[m.cursor].lenFree]
				m.column[m.cursor].row[i] = v
				return int32(m.cursor*capRow) + i
			} else {
				if i := len(m.column[m.cursor].row); i < capRow {
					m.column[m.cursor].row = append(m.column[m.cursor].row, v)
					m.column[m.cursor].indexFree = append(m.column[m.cursor].indexFree, 0)
					return int32(m.cursor*capRow + i)
				}
			}
		}
		m.cursor++
	}

	if m.cursor < capColumn {
		if len(m.indexFree) > 0 {
			m.cursor = int(m.indexFree[0])
			m.indexFree = m.indexFree[1:]
			m.column[m.cursor] = &column[T]{
				row:       []*T{v},
				indexFree: []int32{0},
			}
		} else if m.cursor = len(m.column); m.cursor < capColumn {
			m.column = append(m.column, &column[T]{
				row:       []*T{v},
				indexFree: []int32{0},
			})
		}

		return int32(int(m.cursor) * capRow)
	}

	return -1
}

func (m *Memory[T]) Get(index int32) *T {
	if j := int(index / capRow); j < len(m.column) && m.column[j] != nil {
		if i := int(index % capRow); i < len(m.column[j].row) {
			return m.column[j].row[i]
		}
	}

	return nil
}

func (m *Memory[T]) Free(index int32) {
	if j := int(index / capRow); j < len(m.column) && m.column[j] != nil {
		lenRow := len(m.column[j].row)
		if i := int(index % capRow); i < lenRow && m.column[j].row[i] != nil {
			if lenRow--; i < lenRow {
				m.column[j].row[i] = nil
				m.column[j].indexFree[m.column[j].lenFree] = int32(i)
				m.column[j].lenFree++
			} else {
				m.column[j].row = m.column[j].row[:lenRow]
				m.column[j].indexFree = m.column[j].indexFree[:lenRow]
			}

			if m.column[j].lenFree == len(m.column[j].row) {
				m.column[j] = nil
				for k := 0; k < len(m.indexFree); k++ {
					if j < int(m.indexFree[k]) {
						temp := m.indexFree[k]
						m.indexFree[k] = int16(j)
						j = int(temp)
					}
				}
				m.indexFree = append(m.indexFree, int16(j))

				for len(m.indexFree) > 0 && int(m.indexFree[len(m.indexFree)-1]) == len(m.column)-1 {
					m.column = m.column[:len(m.column)-1]
					m.indexFree = m.indexFree[:len(m.indexFree)-1]
				}
			}
		}
	}
}
