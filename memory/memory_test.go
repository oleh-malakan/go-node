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

func BenchmarkArrayInsert(b *testing.B) {
	b.Run("for", func(b *testing.B) {
		for i := 0; i < b.N; i++ {}
	})

	c := []int{0, 1, 2}
	a := make([]int, 1000000-3)
	a = append(c, a...)
	for i := 3; i < 1000000; i++ {
		a[i] = i + 1
	}
	j := 3
	b.Run("permutations", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			jj := j
			for k := 0; k < len(a); k++ {
				if jj < a[k] {
					temp := a[k]
					a[k] = jj
					jj = temp
				}
			}
			a = append(a, jj)
			j--
		}
	})
	fmt.Println("j: ", j)

	c = []int{0, 1, 2}
	a = make([]int, 1000000-3)
	a = append(c, a...)
	for i := 3; i < 1000000; i++ {
		a[i] = i + 1
	}
	j = 3
	b.Run("slice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b := []int{j}
			b = append(b, a[3:]...)
			a = append(a[:3], b...)		
			j-- 
		}
	})
	fmt.Println("j: ", j)
}