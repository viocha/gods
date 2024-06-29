package main

import (
	"fmt"
	"testing"
)

func Test有序映射(t *testing.T) {
	// 测试TreeMap
	m := NewTreeMap[int, int]()
	for i := range 10 {
		m.Set(i, i)
	}
	fmt.Println("插入0-9之后：")
	m.PrintTree()
	fmt.Println("单行字符串表示：", m)
	fmt.Println()

	fmt.Println("小于等于11：", m.Floor(11).K)
	fmt.Println("小于7：", m.Lower(7).K)
	fmt.Println("大于7：", m.Higher(7).K)
	fmt.Println("大于等于-1：", m.Ceiling(-1).K)
	fmt.Println("前4个节点：")
	m.HeadMap(4).PrintTree()
	fmt.Println("后4个节点：")
	m.TailMap(4).PrintTree()

	fmt.Println("第2个节点：", m.Select(2).K)
	fmt.Println("第15个节点：", m.Select(15))
	fmt.Println("第0个节点：", m.Select(0))
	fmt.Println("15的排名", m.Rank(15))
	fmt.Println("7的排名", m.Rank(7))
	fmt.Println("-1的排名", m.Rank(-1))
	fmt.Println()

	for i := range 5 {
		if m.Has(i) {
			m.Del(i)
		}
	}
	fmt.Println("删除0-4之后的树：")
	m.PrintTree()

	for i := range 10 {
		if m.Has(i) {
			m.Del(i)
		}
	}
	fmt.Println("删除剩余所有节点后：")
	m.PrintTree()

	fmt.Println("=====================带有工厂函数的映射=====================")
	m = NewTreeMap[int, int]().WithFactory(func() int { return 100 })
	fmt.Println("1->", m.Get(1))
	fmt.Println("10->", m.Get(10))
	fmt.Println("大小:", m.Len())
	fmt.Println("字符串表示:", m)
}

func Test有序集(t *testing.T) {
	set := NewTreeSet[int]()
	for i := 9; i >= 0; i-- {
		set.Add(i)
	}
	fmt.Println("插入0-9之后：", set)
	for i := 0; i < 10; i++ {
		set.Del(i)
	}
	fmt.Println("删除0-9之后：", set)
}

func Test多重有序集(t *testing.T) {
	set := NewMultiTreeSet[int]()
	for i := 0; i < 10; i++ {
		set.Add(i)
		set.Add(i)
	}
	fmt.Println("插入0-9之后：", set)
	set.m.PrintTree()
	fmt.Println("每个子树的计数", set.sumMap)
	fmt.Println()

	fmt.Println("5的排名：", set.Rank(5))
	fmt.Println("15的排名：", set.Rank(15))
	fmt.Println("-1的排名：", set.Rank(-1))
	fmt.Println("第5个数：", set.Select(5))
	fmt.Println("第6个数：", set.Select(6))
	fmt.Println("第7个数：", set.Select(7))

	for i := 0; i < 10; i++ {
		set.Del(i)
	}
	fmt.Println("删除0-9之后：", set)
}
