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

	m.Free(t0.index)
	t0.index, _ = m.Put(t0)

	m.Free(t0.index)
	t0.index, _ = m.Put(t0)


	t1 := &test{value: 1}
	t1.index, _ = m.Put(t1)
	t2 := &test{value: 2}
	t2.index, _ = m.Put(t2)

	m.Free(t1.index)
	m.Free(t1.index)
	m.Free(t0.index)
	m.Free(t2.index)
	t0.index, _ = m.Put(t0)
	t1.index, _ = m.Put(t1)
	t2.index, _ = m.Put(t2)

	t3 := &test{value: 3}
	t3.index, _ = m.Put(t3)
	t4 := &test{value: 4}
	t4.index, _ = m.Put(t4)
	t5 := &test{value: 5}
	t5.index, _ = m.Put(t5)
	t6 := &test{value: 6}
	t6.index, _ = m.Put(t6)
	t7 := &test{value: 7}
	t7.index, _ = m.Put(t7)
	t8 := &test{value: 8}
	t8.index, _ = m.Put(t8)
	t9 := &test{value: 9}
	t9.index, _ = m.Put(t9)

	m.Free(t1.index)
	t1.index, _ = m.Put(t1)

	m.Free(t2.index)
	m.Free(t7.index)
	m.Free(t3.index)

	t3.index, _ = m.Put(t3)
	t2.index, _ = m.Put(t2)
	t7.index, _ = m.Put(t7)

	tt := m.Get(t5.index)
	fmt.Println(tt)

	tt = m.Get(55)
	fmt.Println(tt)
}
