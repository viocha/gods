package main

import (
	"fmt"
	"strconv"
	"testing"
)

func Test链表(t *testing.T) {
	l := NewList[int]()
	for i := 0; i < 10; i++ {
		l.PushBack(i)
	}
	fmt.Println(l)
	fmt.Println(l.ToSlice())
	for i := 0; i < 5; i++ {
		l.PopFront()
	}
	fmt.Println(l)
	fmt.Println(l.Len())
}

func Test数组(t *testing.T) {
	a := Range(1, 10)
	a.Ins(3, 10)
	fmt.Println(a)
	fmt.Println("获取最后一个元素：", a.Get(-1))
	fmt.Println("是否包含子数组：", a.HasSubArr(ValArr(10, 4, 5))) // 是否包含子数组
	fmt.Println("求和：", a.Sum())
	fmt.Println("最值：", a.Min(), a.Max())
	fmt.Println("随机值：", a.Choice())
	fmt.Println("映射成为字符串数组：", a.MapToStr(func(x int) string { return "00" + strconv.Itoa(x) }).Join(", "))
	fmt.Println("Filter获取偶数：", a.Filter(func(v int) bool { return v%2 == 0 }))
	fmt.Println("DelFunc删除所有偶数：", a.DelFunc(func(i int, v int) bool { return v%2 == 0 }))

	fmt.Println("ForEach遍历：")
	Arr([]int{1, 2, 3}).ForEach(func(v int) {
		fmt.Println(v)
	})

	arr := ValArr(3, 4, 5)
	fmt.Println("ValArr创建数组：", arr)
	arr = Arr([]int{1, 2, 3})
	fmt.Println("Arr创建数组：", arr)
}

func TestRange(t *testing.T) {
	for i, v := range RangeN(5).GetSlice() {
		fmt.Println(i, v)
	}

	fmt.Println("RangeEq：", RangeEq(1, 10))
	fmt.Println("Range：", Range(1, 10))

	arr2d := ValArr(ValArr(1, 2, 3), ValArr(4, 5, 6))
	fmt.Println("二维数组：", arr2d)
	fmt.Println("二维数组扁平化：", FlatArray(arr2d))

	arr := RangeEq(1, 5)
	fmt.Println("flatMap操作：", arr.FlatMapToInt(func(v int) Array[int] { return RangeEq(1, v) }))
	fmt.Println("flatMap操作：", arr.FlatMap(func(v int) Array[int] { return RangeEq(1, v) }))
}
