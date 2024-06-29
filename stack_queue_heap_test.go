package main

import (
	"fmt"
	"testing"
)

func Test堆(t *testing.T) {
	h := NewHeap[int]().WithMaxCap(6)
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	fmt.Println(h)
	h.Set(4, 100)
	fmt.Println(h.ToSlice())
	h.Del(4)
	fmt.Println(h.ToSlice())
}

func Test栈(t *testing.T) {
	s := NewStack[int]()
	for i := 0; i < 10; i++ {
		s.Push(i)
	}
	fmt.Println(s)
	s.PopUntil(func(x int) bool { return x <= 5 })
	fmt.Println(s)
}

func Test队列(t *testing.T) {
	q := NewQueue[int]()
	for i := 0; i < 10; i++ {
		q.Push(i)
	}
	fmt.Println(q)
	fmt.Println(q.Clone())
	q.PopUntil(func(x int) bool { return x >= 5 })
	fmt.Println(q)
}
