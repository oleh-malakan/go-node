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
	// debug
	for i, column := range m.column {
		// debug
		if column.lenFree > 0 {
			column.lenFree--
			j := column.free[column.lenFree]
			column.row[j] = v
			column.free[column.lenFree] = 0
			return int32(i*rowCap + j), nil
		} else {
			// debug
			if j := len(column.row); j < rowCap {
				column.row = append(column.row, v)
				return int32(i*rowCap + j), nil
			}
		}
	}
	// debug
	if i := len(m.column); i < columnCap {
		m.column = append(m.column, &column[T]{
			row: []*T{v},
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
		if row < len(m.column[column].row) {
			m.column[column].row[row] = nil
			m.column[column].lenFree++
			// debug
			if m.column[column].lenFree < len(m.column[column].free) {
				// debug
				m.column[column].free[m.column[column].lenFree] = row
			} else {
				// debug
				m.column[column].free = append(m.column[column].free, row)
			}
		}
	}
}

type errorEOF struct{}

func (e *errorEOF) Error() string {
	return "EOF"
}
