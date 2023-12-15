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

type continer[T any] struct {
	cont *indexArray[T]
}

type page[T any] struct {
	page *continer[*indexArray[*T]]
}

var EOF = &errorEOF{}

type Memory[T any] struct {
	memory *continer[*page[T]]
}

func (m *Memory[T]) Put(v *T) (int32, error) {

	return 0, EOF
}

func (m *Memory[T]) Get(index int32) *T {
	xy := int(index % zDivider)
	x := xy % cap
	y := xy / cap
	z := int(index / zDivider)
	return m.memory.cont.array[z].page.cont.array[y].array[x]

	return nil
}

func (m *Memory[T]) Free(index int32) {
}

type errorEOF struct{}

func (e *errorEOF) Error() string {
	return "EOF"
}
