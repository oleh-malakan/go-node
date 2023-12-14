package memory

const (
	columnCap = 1000
	rowCap    = 1000000
)

type column[T any] struct {
	row     []*T
	free    []int
	lenFree int
}

var EOF = &errorEOF{}

type Memory[T any] struct {
	column []*column[T]
}

func (m *Memory[T]) Put(v *T) (int32, error) {
	for i, column := range m.column {
		if column.lenFree > 0 {
			column.lenFree--
			j := column.free[column.lenFree]
			column.row[j] = v
			return int32(i*rowCap + j), nil
		} else {
			if j := len(column.row); j < rowCap {
				column.row = append(column.row, v)
				column.free = append(column.free, 0)
				return int32(i*rowCap + j), nil
			}
		}
	}
	if i := len(m.column); i < columnCap {
		m.column = append(m.column, &column[T]{
			row:  []*T{v},
			free: []int{0},
		})
		return int32(i * rowCap), nil
	}

	return 0, EOF
}

func (m *Memory[T]) Get(index int32) *T {
	column := int(index / rowCap)
	row := int(index % rowCap)
	if column < len(m.column) {
		if row < len(m.column[column].row) {
			return m.column[column].row[row]
		}
	}

	return nil
}

func (m *Memory[T]) Free(index int32) {
	column := int(index / rowCap)
	row := int(index % rowCap)
	if column < len(m.column) {
		if row < len(m.column[column].row) && m.column[column].row[row] != nil {
			m.column[column].row[row] = nil
			m.column[column].free[m.column[column].lenFree] = row
			m.column[column].lenFree++
			if m.column[column].lenFree == len(m.column[column].row) {
				m.column[column].row = nil
				m.column[column].free = nil
				m.column[column].lenFree = 0
			}
		}
	}
}

type errorEOF struct{}

func (e *errorEOF) Error() string {
	return "EOF"
}
