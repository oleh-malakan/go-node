package memory

import (
	"fmt"
	"testing"
)

func TestMemory(t *testing.T) {
	index := 9999999
	row := index % rowCap
	column := index / rowCap
	index = column*rowCap + row
	fmt.Println(row, column, index)

	type test struct {
		value int
		index int32
	}

	m := &Memory[test]{}

	var err error
	t0 := &test{value: 0}
	t0.index, err = m.Put(t0)
	if err != nil {
		fmt.Println(err)
	}
	t1 := &test{value: 1}
	t1.index, _ = m.Put(t1)
	t2 := &test{value: 2}
	t2.index, _ = m.Put(t2)

	m.Free(t1.index)
	t1.index, _ = m.Put(t1)

	tt := m.Get(55)
	fmt.Println(tt)

}
