package memory

import (
	"fmt"
	"testing"
)

type test struct {
	value int
	index int32
}

func TestMemory(t *testing.T) {
	m := &Memory[test]{}

	t0 := &test{value: 0}
	t0.index = m.Put(t0)

	m.Free(t0.index)
	t0.index = m.Put(t0)

	m.Free(t0.index)
	t0.index = m.Put(t0)

	t1 := &test{value: 1}
	t1.index = m.Put(t1)
	t2 := &test{value: 2}
	t2.index = m.Put(t2)

	m.Free(t1.index)
	m.Free(t1.index)
	m.Free(t0.index)
	m.Free(t2.index)
	t0.index = m.Put(t0)
	t1.index = m.Put(t1)
	t2.index = m.Put(t2)

	t3 := &test{value: 3}
	t3.index = m.Put(t3)
	t4 := &test{value: 4}
	t4.index = m.Put(t4)
	t5 := &test{value: 5}
	t5.index = m.Put(t5)
	t6 := &test{value: 6}
	t6.index = m.Put(t6)
	t7 := &test{value: 7}
	t7.index = m.Put(t7)
	t8 := &test{value: 8}
	t8.index = m.Put(t8)
	t9 := &test{value: 9}
	t9.index = m.Put(t9)

	m.Free(t1.index)
	t1.index = m.Put(t1)

	m.Free(t2.index)
	m.Free(t7.index)
	m.Free(t3.index)

	m.Free(t9.index)
	m.Free(t8.index)

	t3.index = m.Put(t3)
	t2.index = m.Put(t2)
	t7.index = m.Put(t7)

	tt := m.Get(t5.index)
	fmt.Println(tt)

	tt = m.Get(55)
	fmt.Println(tt)

	t9.index = m.Put(t9)
	t8.index = m.Put(t8)

	t10 := &test{value: 10}
	t10.index = m.Put(t10)

	m.Free(t10.index)
	t10.index = m.Put(t10)

	m.Free(t10.index)
	m.Free(t9.index)

	t10.index = m.Put(t10)
	t9.index = m.Put(t9)

	m.Free(t10.index)
	m.Free(t9.index)

	t9.index = m.Put(t9)
	t10.index = m.Put(t10)

	m.Free(t0.index)
	m.Free(t1.index)
	m.Free(t2.index)
	m.Free(t3.index)
	m.Free(t4.index)
	m.Free(t5.index)
	m.Free(t6.index)
	m.Free(t7.index)
	m.Free(t8.index)
	m.Free(t9.index)

	t9.index = m.Put(t9)

	m.Free(t9.index)
	m.Free(t10.index)

	t9.index = m.Put(t9)

}
