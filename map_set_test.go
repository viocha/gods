package gods

import (
	"fmt"
	"testing"
)

func Test哈希表(t *testing.T) {
	m := NewHashMap[int, int]()
	for i := 0; i < 10; i++ {
		m.Set(i, i)
	}
	fmt.Println(m)
}

func Test哈希集(t *testing.T) {
	s := NewHashSet[int]()
	for i := 0; i < 10; i++ {
		s.Add(i)
	}
	fmt.Println(s)
}

func Test链式哈希表(t *testing.T) {
	m := NewLinkedHashMap[int, int]().WithAccessOrderMode().WithMaxCap(10)
	for i := 15; i > 0; i-- {
		m.Set(i, i)
		if i == 7 {
			m.Get(12)
			m.Get(14)
		}
	}
	fmt.Println(m)
	internalMap := make(map[int]int, m.Len())
	for k, v := range m.m {
		internalMap[k] = v.v
	}
	fmt.Println(internalMap)
}

func Test链式哈希集(t *testing.T) {
	s := NewLinkedHashSet[int]().WithAccessOrderMode().WithMaxCap(10)
	for i := 15; i > 0; i-- {
		s.Add(i)
		if i == 7 {
			s.Has(12)
			s.Has(14)
		}
	}
	fmt.Println(s.m)
	fmt.Println(s)
	fmt.Println(s.Len())

	fmt.Println("第一个元素：", s.First())
	fmt.Println("最后一个元素：", s.Last())
	fmt.Println("7的下一个元素：", s.Next(7))
	fmt.Println("6的前一个元素：", s.Prev(6))
}

func Test多重哈希集(t *testing.T) {
	s := NewMultiHashSet[int]()
	for i := 0; i < 10; i++ {
		s.Add(i).Add(i)
	}
	fmt.Println(s)
	fmt.Println(s.ToHashSet())
	fmt.Println(s.ToMap())
	for i := 0; i < 10; i++ {
		s.Del(i)
	}
	fmt.Println(s)
}
