package memory

import (
	"fmt"
	"testing"
)

func TestMemory(t *testing.T) {
	index := 77871032
	xy := index % depth
	x := xy % cap
	y := xy / cap
	z := index / depth
	index = int(z)*depth + int(y)*cap + int(x)
	fmt.Println(x, y, z, index)

	type test struct {
		value int
		index int16
	}

	a := &indexArray[test]{}

	var err error
	t0 := &test{value: 0}
	t0.index = a.put(t0)
	if err != nil {
		fmt.Println(err)
	}

	a.free(t0.index)
	t0.index = a.put(t0)

	a.free(t0.index)
	t0.index = a.put(t0)

	t1 := &test{value: 1}
	t1.index = a.put(t1)
	t2 := &test{value: 2}
	t2.index = a.put(t2)

	a.free(t1.index)
	a.free(t1.index)
	a.free(t0.index)
	a.free(t2.index)
	t0.index = a.put(t0)
	t1.index = a.put(t1)
	t2.index = a.put(t2)

	t3 := &test{value: 3}
	t3.index = a.put(t3)
	t4 := &test{value: 4}
	t4.index = a.put(t4)
	t5 := &test{value: 5}
	t5.index = a.put(t5)
	t6 := &test{value: 6}
	t6.index = a.put(t6)
	t7 := &test{value: 7}
	t7.index = a.put(t7)
	t8 := &test{value: 8}
	t8.index = a.put(t8)
	t9 := &test{value: 9}
	t9.index = a.put(t9)

	a.free(t1.index)
	t1.index = a.put(t1)

	a.free(t2.index)
	a.free(t7.index)
	a.free(t3.index)

	a.free(t9.index)
	a.free(t8.index)

	t3.index = a.put(t3)
	t2.index = a.put(t2)
	t7.index = a.put(t7)

	tt := a.get(t5.index)
	fmt.Println(tt)

	tt = a.get(55)
	fmt.Println(tt)
}
