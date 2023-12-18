package memory

import (
	"fmt"
	"testing"
)

type test struct {
	value int
	index int
}

func TestMain(t *testing.T) {
	index := 77871032
	xy := index % depth
	x := xy % cap
	y := xy / cap
	z := index / depth
	index = int(z)*depth + int(y)*cap + int(x)
	fmt.Println(x, y, z, index)
}

func TestColumn(t *testing.T) {
	c := &column[test]{}

	t0 := &test{value: 0}
	t0.index = c.put(t0)

	c.free(t0.index)
	t0.index = c.put(t0)

	c.free(t0.index)
	t0.index = c.put(t0)

	t1 := &test{value: 1}
	t1.index = c.put(t1)
	t2 := &test{value: 2}
	t2.index = c.put(t2)

	c.free(t1.index)
	c.free(t1.index)
	c.free(t0.index)
	c.free(t2.index)
	t0.index = c.put(t0)
	t1.index = c.put(t1)
	t2.index = c.put(t2)

	t3 := &test{value: 3}
	t3.index = c.put(t3)
	t4 := &test{value: 4}
	t4.index = c.put(t4)
	t5 := &test{value: 5}
	t5.index = c.put(t5)
	t6 := &test{value: 6}
	t6.index = c.put(t6)
	t7 := &test{value: 7}
	t7.index = c.put(t7)
	t8 := &test{value: 8}
	t8.index = c.put(t8)
	t9 := &test{value: 9}
	t9.index = c.put(t9)

	c.free(t1.index)
	t1.index = c.put(t1)

	c.free(t2.index)
	c.free(t7.index)
	c.free(t3.index)

	c.free(t9.index)
	c.free(t8.index)

	t3.index = c.put(t3)
	t2.index = c.put(t2)
	t7.index = c.put(t7)

	tt := c.get(t5.index)
	fmt.Println(tt)

	tt = c.get(55)
	fmt.Println(tt)
}

func TestOrderVector(t *testing.T) {
	v := *&orderVector[test]{}

	t0 := &test{value: 0}
	t0.index = v.put(t0)

	v.free(t0.index)
	t0.index = v.put(t0)

	v.free(t0.index)
	t0.index = v.put(t0)

	t1 := &test{value: 1}
	t1.index = v.put(t1)
	t2 := &test{value: 2}
	t2.index = v.put(t2)

	v.free(t1.index)
	v.free(t1.index)
	v.free(t0.index)
	v.free(t2.index)
	t0.index = v.put(t0)
	t1.index = v.put(t1)
	t2.index = v.put(t2)

	t3 := &test{value: 3}
	t3.index = v.put(t3)
	t4 := &test{value: 4}
	t4.index = v.put(t4)
	t5 := &test{value: 5}
	t5.index = v.put(t5)
	t6 := &test{value: 6}
	t6.index = v.put(t6)
	t7 := &test{value: 7}
	t7.index = v.put(t7)
	t8 := &test{value: 8}
	t8.index = v.put(t8)
	t9 := &test{value: 9}
	t9.index = v.put(t9)

	v.free(t1.index)
	t1.index = v.put(t1)

	v.free(t2.index)
	v.free(t7.index)
	v.free(t3.index)

	v.free(t9.index)
	v.free(t8.index)

	t3.index = v.put(t3)
	t2.index = v.put(t2)
	t7.index = v.put(t7)

	tt := v.get(t5.index)
	fmt.Println(tt)

	tt = v.get(55)
	fmt.Println(tt)
}