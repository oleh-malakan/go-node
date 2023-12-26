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
	column     []*column[T]
	indexFree  []int
	columnFree *column[T]
	cursor     int
	len        int
}

func (m *Memory[T]) Put(v *T) int32 {
	if v != nil {
		m.len++
		for m.cursor < len(m.column) {
			if m.column[m.cursor] != nil && m.column[m.cursor].lenFree > 0 {
				m.column[m.cursor].lenFree--
				i := m.column[m.cursor].indexFree[m.column[m.cursor].lenFree]
				m.column[m.cursor].row[i] = v
				return int32(m.cursor*capRow) + i
			}
			m.cursor++
		}

		if len(m.indexFree) > 0 {
			m.cursor = int(m.indexFree[0])
			m.indexFree = m.indexFree[1:]
			if m.columnFree == nil {
				m.column[m.cursor] = &column[T]{
					row:       make([]*T, capRow),
					indexFree: make([]int32, capRow),
				}
				for i := 0; i < capRow; i++ {
					m.column[m.cursor].indexFree[i] = int32(capRow - i - 1)
				}
				m.column[m.cursor].lenFree = capRow
			} else {
				m.column[m.cursor] = m.columnFree
				m.columnFree = nil
			}
			m.column[m.cursor].lenFree--
			i := m.column[m.cursor].indexFree[m.column[m.cursor].lenFree]
			m.column[m.cursor].row[i] = v
			return int32(m.cursor*capRow) + i
		} else if m.cursor < capColumn {
			if m.columnFree == nil {
				m.column = append(m.column, &column[T]{
					row:       make([]*T, capRow),
					indexFree: make([]int32, capRow),
				})
				for i := 0; i < capRow; i++ {
					m.column[m.cursor].indexFree[i] = int32(capRow - i - 1)
				}
				m.column[m.cursor].lenFree = capRow
			} else {
				m.column = append(m.column, m.columnFree)
				m.columnFree = nil
			}
			m.column[m.cursor].lenFree--
			i := m.column[m.cursor].indexFree[m.column[m.cursor].lenFree]
			m.column[m.cursor].row[i] = v
			return int32(m.cursor*capRow) + i
		}

		m.len--
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
		if i := int(index % capRow); i < capRow && m.column[j].row[i] != nil {
			if j < m.cursor {
				m.cursor = j
			}
			m.len--
			m.column[j].row[i] = nil
			m.column[j].indexFree[m.column[j].lenFree] = int32(i)
			m.column[j].lenFree++

			if m.column[j].lenFree == len(m.column[j].row) {
				m.columnFree = m.column[j]
				m.column[j] = nil

				for k := 0; k < len(m.indexFree); k++ {
					if j < int(m.indexFree[k]) {
						temp := m.indexFree[k]
						m.indexFree[k] = j
						j = int(temp)
					}
				}
				m.indexFree = append(m.indexFree, j)

				for len(m.indexFree) > 0 && int(m.indexFree[len(m.indexFree)-1]) == len(m.column)-1 {
					m.column = m.column[:len(m.column)-1]
					m.indexFree = m.indexFree[:len(m.indexFree)-1]
				}
			}
		}
	}
}

func (m *Memory[T]) Len() int {
	return m.len
}
