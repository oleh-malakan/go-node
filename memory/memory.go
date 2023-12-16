package memory

const (
	cap      = 1000
	zDivider = 1000000
)

type indexArray[T any] struct {
	array     []T
	indexFree []int16
	lenFree   int16
}

func (a *indexArray[T]) put(v *T) int16 {
	return 0
}

func (a *indexArray[T]) get(index int16) *T {
	return nil
}

func (a *indexArray[T]) free(index int16) {

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

type memory[T any] struct {
	m *orderContainer[*page[T]]
}

func (m *memory[T]) open() (int16, *page[T]) {
	return 0, nil
}

func (m *memory[T]) get(index int16) *page[T] {
	return nil
}

func (m *memory[T]) free(index int16) {

}

var EOF = &errorEOF{}

type Memory[T any] struct {
	memory *memory[T]
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
