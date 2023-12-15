package memory

const (
	cap      = 1000
	zDivider = 1000000
)

type indexArray[T any] struct {
	array   []T
	free    []int32
	lenFree int32
}

type page[T any] struct {
	page *indexArray[*indexArray[*T]]
}

var EOF = &errorEOF{}

type Memory[T any] struct {
	memory *indexArray[*page[T]]
}

func (m *Memory[T]) Put(v *T) (int32, error) {

	return 0, EOF
}

func (m *Memory[T]) Get(index int32) *T {
	xy := int(index % zDivider)
	x := xy % cap
	y := xy / cap
	z := int(index / zDivider)

	return nil
}

func (m *Memory[T]) Free(index int32) {
}

type errorEOF struct{}

func (e *errorEOF) Error() string {
	return "EOF"
}
