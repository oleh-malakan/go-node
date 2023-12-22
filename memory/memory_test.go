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
	c := &Memory[test]{}

	t0 := &test{value: 0}
	t0.index = c.Put(t0)

	c.Free(t0.index)
	t0.index = c.Put(t0)

	c.Free(t0.index)
	t0.index = c.Put(t0)

	t1 := &test{value: 1}
	t1.index = c.Put(t1)
	t2 := &test{value: 2}
	t2.index = c.Put(t2)

	c.Free(t1.index)
	c.Free(t1.index)
	c.Free(t0.index)
	c.Free(t2.index)
	t0.index = c.Put(t0)
	t1.index = c.Put(t1)
	t2.index = c.Put(t2)

	t3 := &test{value: 3}
	t3.index = c.Put(t3)
	t4 := &test{value: 4}
	t4.index = c.Put(t4)
	t5 := &test{value: 5}
	t5.index = c.Put(t5)
	t6 := &test{value: 6}
	t6.index = c.Put(t6)
	t7 := &test{value: 7}
	t7.index = c.Put(t7)
	t8 := &test{value: 8}
	t8.index = c.Put(t8)
	t9 := &test{value: 9}
	t9.index = c.Put(t9)

	c.Free(t1.index)
	t1.index = c.Put(t1)

	c.Free(t2.index)
	c.Free(t7.index)
	c.Free(t3.index)

	c.Free(t9.index)
	c.Free(t8.index)

	t3.index = c.Put(t3)
	t2.index = c.Put(t2)
	t7.index = c.Put(t7)

	tt := c.Get(t5.index)
	fmt.Println(tt)

	tt = c.Get(55)
	fmt.Println(tt)
}
