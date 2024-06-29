package gods

import (
	"cmp"
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// ===================================双向链表===================================

// --------------------双向链表节点--------------------

type ListNode[T any] struct {
	Val        T
	prev, next *ListNode[T]
	list       *List[T]
}

func newListNode[T any](val T, list *List[T]) *ListNode[T] { return &ListNode[T]{Val: val, list: list} }
func (o *ListNode[T]) insAfter(mark *ListNode[T]) *ListNode[T] {
	nex := mark.next
	mark.next, nex.prev = o, o
	o.prev, o.next = mark, nex
	*o.list.len++
	return o
}
func (o *ListNode[T]) insBefore(mark *ListNode[T]) *ListNode[T] {
	pre := mark.prev
	pre.next, mark.prev = o, o
	o.prev, o.next = pre, mark
	*o.list.len++
	return o
}
func (o *ListNode[T]) remove() *ListNode[T] {
	pre, nex := o.prev, o.next
	pre.next, nex.prev = nex, pre
	o.prev, o.next = nil, nil
	*o.list.len--
	return o
}

func (o *ListNode[T]) Next() *ListNode[T] {
	if o.next == o.list.root {
		return nil
	}
	return o.next
}
func (o *ListNode[T]) Prev() *ListNode[T] {
	if o.prev == o.list.root {
		return nil
	}
	return o.prev
}

// --------------------链表定义--------------------

type List[T any] struct {
	root *ListNode[T]
	len  *int
}

func NewList[T any]() List[T] {
	o := List[T]{len: new(int)}
	root := o.NewNode(*new(T))
	root.next, root.prev = root, root
	o.root = root
	return o
}
func (o List[T]) NewNode(val T) *ListNode[T] { return newListNode[T](val, &o) }

// ---------------ValueContainer接口---------------
func (o List[T]) Len() int { return *o.len }
func (o List[T]) Clear() List[T] {
	root := o.root
	root.next, root.prev = root, root
	*o.len = 0
	return o
}
func (o List[T]) Clone() List[T] {
	res := NewList[T]()
	o.ForEach(func(v T) { res.PushBack(v) })
	return res
}
func (o List[T]) ForEach(f func(T)) List[T] {
	vals := make([]T, 0, o.Len())
	for p := o.Front(); p != nil; p = p.Next() {
		vals = append(vals, p.Val)
	}
	for _, v := range vals { // 支持同时遍历和修改
		f(v)
	}
	return o
}
func (o List[T]) ToSlice() []T {
	sl := make([]T, 0, o.Len())
	o.ForEach(func(v T) { sl = append(sl, v) })
	return sl
}
func (o List[T]) String() string { return "List" + fmt.Sprint(o.ToSlice()) }

// --------------------List接口--------------------

// 基本操作
func (o List[T]) Front() *ListNode[T] {
	if o.Len() == 0 {
		return nil
	}
	return o.root.next
}
func (o List[T]) Back() *ListNode[T] {
	if o.Len() == 0 {
		return nil
	}
	return o.root.prev
}

func (o List[T]) PushFront(val T) *ListNode[T] { return o.InsAfter(val, o.root) }
func (o List[T]) PushBack(val T) *ListNode[T]  { return o.InsBefore(val, o.root) }
func (o List[T]) InsAfter(val T, mark *ListNode[T]) *ListNode[T] {
	return o.NewNode(val).insAfter(mark)
}
func (o List[T]) InsBefore(val T, mark *ListNode[T]) *ListNode[T] {
	return o.NewNode(val).insBefore(mark)
}

func (o List[T]) Remove(node *ListNode[T]) T { return node.remove().Val }
func (o List[T]) PopFront() *ListNode[T]     { return o.root.next.remove() }
func (o List[T]) PopBack() *ListNode[T]      { return o.root.prev.remove() }

// 移动节点
func (o List[T]) MoveAfter(node, mark *ListNode[T]) T  { return node.remove().insAfter(mark).Val }
func (o List[T]) MoveBefore(node, mark *ListNode[T]) T { return node.remove().insBefore(mark).Val }
func (o List[T]) MoveToFront(node *ListNode[T]) T      { return o.MoveAfter(node, o.root) }
func (o List[T]) MoveToBack(node *ListNode[T]) T       { return o.MoveBefore(node, o.root) }

// ===================================动态数组===================================
type Array[T any] struct {
	data *[]T
	cmp  func(T, T) int
	eq   func(T, T) bool
}

// 创建空数组
func NewArray[T any]() Array[T] {
	var compare func(T, T) int
	Val := reflect.ValueOf
	v := Val(*new(T))

	if v.CanInt() {
		compare = func(a, b T) int {
			return cmp.Compare[int64](Val(a).Int(), Val(b).Int())
		}
	} else if v.CanUint() {
		compare = func(a, b T) int {
			return cmp.Compare[uint64](Val(a).Uint(), Val(b).Uint())
		}
	} else if v.CanFloat() {
		compare = func(a, b T) int {
			return cmp.Compare[float64](Val(a).Float(), Val(b).Float())
		}
	} else if v.Kind() == reflect.String {
		compare = func(a, b T) int {
			return cmp.Compare[string](Val(a).String(), Val(b).String())
		}
	}

	eq := func(a, b T) bool { return Val(a).Equal(Val(b)) }

	o := Array[T]{data: new([]T), cmp: compare, eq: eq}
	return o
}

// 创建一个具有大小n的数组
func MakeArray[T any](n int) Array[T] {
	arr := NewArray[T]()
	data := make([]T, n)
	arr.data = &data
	return arr
}

// 复制切片的内容，创建数组
func NewArrayFromSlice[T any](data []T) Array[T] {
	c := slices.Clone(data)
	return NewArray[T]().WithSlice(c)
}

// 引用切片构建数组
func Arr[T any](data []T) Array[T] { return NewArray[T]().WithSlice(data) }

// 根据数组值构建数组
func ValArr[T any](v ...T) Array[T] { return NewArray[T]().WithSlice(v) }

// 修改内部的切片
func (o Array[T]) WithSlice(data []T) Array[T] { o.data = &data; return o }

// 获取内部的切片
func (o Array[T]) GetSlice() []T { return *o.data }

// 同时获取内部切片和长度
func (o Array[T]) getData() ([]T, int) { return *o.data, len(*o.data) }

// --------------------ValueContainer接口--------------------
func (o Array[T]) Len() int                   { return len(*o.data) }
func (o Array[T]) Clone() Array[T]            { return NewArrayFromSlice(*o.data) }
func (o Array[T]) Clear() Array[T]            { *o.data = (*o.data)[:0]; return o }
func (o Array[T]) ForEach(f func(T)) Array[T] { return o.ForEachIdxVal(func(i int, v T) { f(v) }) }
func (o Array[T]) ToSlice() []T               { return slices.Clone(*o.data) }
func (o Array[T]) String() string             { return "Array" + fmt.Sprint(*o.data) }

// --------------------Array接口--------------------

// ---------------遍历和转换---------------

func (o Array[T]) Keys() Array[int]                { return RangeN(o.Len()) }
func (o Array[T]) ForEachIdx(f func(int)) Array[T] { return o.ForEachIdxVal(func(i int, v T) { f(i) }) }
func (o Array[T]) ForEachIdxVal(f func(int, T)) Array[T] {
	for i, v := range o.ToSlice() {
		f(i, v)
	}
	return o
}

// ---------------基本操作---------------
// 查询值，支持负值索引
func (o Array[T]) Get(i int) T             { return (*o.data)[o.idx(i)] }
func (o Array[T]) Set(i int, v T) Array[T] { (*o.data)[o.idx(i)] = v; return o }

// 截取子数组视图
func (o Array[T]) Slice(i, j int) Array[T] {
	i, j = o.idx(i), o.idx(j)
	*o.data = (*o.data)[i:j]
	return o
}

// 截取子数组副本
func (o Array[T]) Sliced(i, j int) Array[T] {
	i, j = o.idx(i), o.idx(j)
	return NewArray[T]().WithSlice((*o.data)[i:j])
}

// 插入和删除元素
func (o Array[T]) Ins(i int, v T) Array[T] { return o.Replace(i, i, v) }
func (o Array[T]) Del(i int) Array[T]      { return o.Replace(i, i+1) }
func (o Array[T]) Push(v ...T) Array[T] {
	*o.data = append(*o.data, v...)
	return o
}
func (o Array[T]) Pop() T {
	x := o.Get(-1)
	o.Del(-1)
	return x
}

// 支持负值索引
func (o Array[T]) idx(i int) int {
	preIdx := i
	if i < 0 {
		i = len(*o.data) + i
	}
	if i < 0 || i >= len(*o.data) {
		panic("索引超出范围：" + strconv.Itoa(preIdx))
	}
	return i
}

// ---------------组合操作---------------

// 将i..j之间的元素替换成值列表
func (o Array[T]) Replace(i, j int, v ...T) Array[T] {
	i, j = o.idx(i), o.idx(j)
	*o.data = slices.Replace(*o.data, i, j, v...)
	return o
}

// 替换每个元素
func (o Array[T]) ReplaceFunc(f func(int, T) T) Array[T] {
	return o.ForEachIdxVal(func(i int, v T) {
		o.Set(i, f(i, v))
	})
}

// 如果f返回true，则删除该元素
func (o Array[T]) DelFunc(f func(int, T) bool) Array[T] {
	var newData []T
	o.ForEachIdxVal(func(i int, v T) {
		if !f(i, v) {
			newData = append(newData, v)
		}
	})
	*o.data = newData
	return o
}

// 将另一个数组到当前数组中
func (o Array[T]) Extend(other Array[T]) Array[T] { return o.Push(*other.data...) }

// 填充数组每个元素
func (o Array[T]) Fill(x T) Array[T] {
	return o.ForEachIdx(func(i int) {
		o.Set(i, x)
	})
}

func (o Array[T]) FillFunc(f func(int) T) Array[T] {
	return o.ForEachIdx(func(i int) {
		o.Set(i, f(i))
	})
}

// 原地反转
func (o Array[T]) Reverse() Array[T]  { slices.Reverse(*o.data); return o }
func (o Array[T]) Reversed() Array[T] { return o.Clone().Reverse() }

// ---------------随机化操作---------------

func (o Array[T]) Shuffle() Array[T] {
	d, n := o.getData()
	rand.Shuffle(n, func(i, j int) { d[i], d[j] = d[j], d[i] })
	return o
}
func (o Array[T]) Choice() T {
	i := rand.Intn(o.Len())
	return o.Get(i)
}

// ---------------查找，排序，比较---------------

// 判断是否包含元素
func (o Array[T]) Has(x T) bool                 { return o.Index(x) >= 0 }
func (o Array[T]) HasFunc(eq func(T) bool) bool { return o.IndexFunc(eq) >= 0 }
func (o Array[T]) HasBi(x T) bool               { return o.IndexBi(x) >= 0 }

// 二分查找是否包含元素
func (o Array[T]) HasBiFunc(x T, cmp func(T, T) int) bool { return o.IndexBiFunc(x, cmp) >= 0 }

// 线性查找和元素相等的索引，不存在返回-1
func (o Array[T]) Index(x T) int                 { return o.IndexFunc(func(v T) bool { return o.eq(x, v) }) }
func (o Array[T]) IndexFunc(eq func(T) bool) int { return slices.IndexFunc(*o.data, eq) }

// 二分查找和元素相等的索引，不存在返回-1
func (o Array[T]) IndexBi(x T) int { return o.IndexBiFunc(x, o.cmp) }
func (o Array[T]) IndexBiFunc(x T, cmp func(T, T) int) int {
	i, found := slices.BinarySearchFunc(*o.data, x, cmp)
	if found {
		return i
	}
	return -1
}

// 二分查找第一个>=x的索引，不存则返回n
func (o Array[T]) BiSearch(x T) int {
	return o.BiSearchFunc(func(v T) bool { return o.cmp(v, x) >= 0 })
}

// 二分查找第一个使check函数返回true的索引，会假定f在数组左段返回false，右段返回true
func (o Array[T]) BiSearchFunc(check func(T) bool) int {
	d, n := o.getData()
	return sort.Search(n, func(i int) bool { return check(d[i]) })
}

// 返回第一个使check函数返回true的元素值，如果不存在返回类型零值
func (o Array[T]) Find(check func(T) bool) T {
	i := o.IndexFunc(check)
	if 0 <= i && i < o.Len() {
		return o.Get(i)
	}
	return *new(T)
}

// 判断是否包含子数组
func (o Array[T]) HasSubArr(sub Array[T]) bool { return o.IndexSubArr(sub) >= 0 }
func (o Array[T]) HasSubArrFunc(sub Array[T], eq func(T, T) bool) bool {
	return o.IndexSubArrFunc(sub, eq) >= 0
}
func (o Array[T]) IndexSubArr(sub Array[T]) int { return o.IndexSubArrFunc(sub, o.eq) }

// 查找子数组位置（KMP算法），使用自定义比较函数
func (o Array[T]) IndexSubArrFunc(sub Array[T], eq func(T, T) bool) int {
	// 构造next数组
	next := make([]int, sub.Len())
	j := 0
	for i := 1; i < len(next); i++ {
		for j > 0 && !eq(sub.Get(i), sub.Get(j)) {
			j = next[j-1]
		}
		if eq(sub.Get(i), sub.Get(j)) {
			j++
		}
		next[i] = j
	}

	i := 0
	j = 0
	for i < o.Len() && j < sub.Len() {
		if eq(o.Get(i), sub.Get(j)) {
			i, j = i+1, j+1
		} else if j > 0 {
			j = next[j-1]
		} else {
			i++
		}
	}
	if j == sub.Len() {
		return i - j
	}
	return -1
}

// 稳定排序
func (o Array[T]) Sort() Array[T] { return o.SortFunc(o.cmp) } // 默认稳定排序
func (o Array[T]) SortFunc(cmp func(T, T) int) Array[T] {
	slices.SortStableFunc(*o.data, cmp)
	return o
}

// 按元素比较顺序
func (o Array[T]) Cmp(other Array[T]) int { return o.CmpFunc(other, o.cmp) }
func (o Array[T]) CmpFunc(other Array[T], f func(T, T) int) int {
	return slices.CompareFunc(*o.data, *other.data, f)
}

// 判断所有元素是否相等
func (o Array[T]) Equal(other Array[T]) bool { return o.EqualFunc(other, o.eq) }
func (o Array[T]) EqualFunc(other Array[T], f func(T, T) bool) bool {
	return slices.EqualFunc(*o.data, *other.data, f)
}

// ---------------聚合操作---------------

func (o Array[T]) Min() T                     { return o.MinFunc(o.cmp) }
func (o Array[T]) MinFunc(f func(T, T) int) T { return slices.MinFunc(*o.data, f) }
func (o Array[T]) Max() T                     { return o.MaxFunc(o.cmp) }
func (o Array[T]) MaxFunc(f func(T, T) int) T { return slices.MaxFunc(*o.data, f) }

// 转换成int,uint或者float64，然后求和
func (o Array[T]) Sum() T {
	Val := reflect.ValueOf
	res := new(T)
	v := reflect.ValueOf(res).Elem()
	if v.CanInt() {
		var sum int64
		for _, x := range *o.data {
			sum += Val(x).Int()
		}
		v.SetInt(sum)
	} else if v.CanUint() {
		var sum uint64
		for _, x := range *o.data {
			sum += Val(x).Uint()
		}
		v.SetUint(sum)
	} else if v.CanFloat() {
		var sum float64
		for _, x := range *o.data {
			sum += Val(x).Float()
		}
		v.SetFloat(sum)
	} else {
		panic("类型无法进行求和：" + v.Type().String())
	}
	return *res
}

// 如果通过判断是否和零值相等，转换成true或者false
func (o Array[T]) Any() bool {
	zero := *new(T)
	for _, x := range *o.data {
		if o.cmp(x, zero) != 0 {
			return true
		}
	}
	return false
}
func (o Array[T]) All() bool {
	zero := *new(T)
	for _, x := range *o.data {
		if o.cmp(x, zero) == 0 {
			return false
		}
	}
	return true
}

// 使用fmt.Sprint转换成字符串，然后连接
func (o Array[T]) Join(seq string) string {
	return strings.Join(o.MapToStr(func(x T) string { return fmt.Sprint(x) }).GetSlice(), seq)
}

// ---------------函数式操作(返回一个副本)---------------
func (o Array[T]) Filter(f func(T) bool) Array[T] {
	newData := make([]T, 0)
	o.ForEachIdxVal(func(i int, v T) {
		if f(v) {
			newData = append(newData, v)
		}
	})
	return Arr(newData)
}

func (o Array[T]) Reduce(f func(T, T) T) T {
	if o.Len() == 0 {
		return *new(T)
	}
	res := o.Get(0)
	for i := 1; i < o.Len(); i++ {
		res = f(res, o.Get(i))
	}
	return res
}
func (o Array[T]) ReduceInitial(f func(T, T) T, initial T) T {
	if o.Len() == 0 {
		return initial
	}
	res := initial
	o.ForEach(func(v T) { res = f(res, v) })
	return res
}
func (o Array[T]) Map(f func(T) T) Array[T] {
	o = o.Clone()
	return o.ForEachIdxVal(func(i int, v T) {
		o.Set(i, f(v))
	})
}
func (o Array[T]) MapToInt(f func(T) int) Array[int] {
	return MapArrayTo(o, func(x T) int { return f(x) })
}
func (o Array[T]) MapToFloat(f func(T) float64) Array[float64] {
	return MapArrayTo(o, func(x T) float64 { return f(x) })
}
func (o Array[T]) MapToStr(f func(T) string) Array[string] {
	return MapArrayTo(o, func(x T) string { return f(x) })
}
func (o Array[T]) MapToBool(f func(T) bool) Array[bool] {
	return MapArrayTo(o, func(x T) bool { return f(x) })
}

// 扁平化map操作
func (o Array[T]) FlatMap(f func(T) Array[T]) Array[T] {
	//return FlatMapArrayTo(o, f) // TODO: 出现错误 instantiation cycle
	res := NewArray[T]()
	o.ForEach(func(v T) {
		res.Push(f(v).GetSlice()...)
	})
	return res
}
func (o Array[T]) FlatMapToInt(f func(T) Array[int]) Array[int] { return FlatMapArrayTo(o, f) }
func (o Array[T]) FlatMapToFloat(f func(T) Array[float64]) Array[float64] {
	return FlatMapArrayTo(o, f)
}
func (o Array[T]) FlatMapToStr(f func(T) Array[string]) Array[string] {
	return FlatMapArrayTo(o, f)
}
func (o Array[T]) FlatMapToBool(f func(T) Array[bool]) Array[bool] {
	return FlatMapArrayTo(o, f)
}

// 支持转换成任意类型的map操作
func MapArrayTo[T, R any](arr Array[T], f func(T) R) Array[R] {
	res := NewArray[R]()
	arr.ForEach(func(v T) { res.Push(f(v)) })
	return res
}

func FlatMapArrayTo[T, R any](array Array[T], f func(T) Array[R]) Array[R] {
	res := NewArray[Array[R]]()
	array.ForEach(func(v T) { res.Push(f(v)) })
	return FlatArray(res)
}

// --------------------其他操作--------------------

// 将数组扁平化，支持Array[Array[T]]
func FlatArray[T any](array Array[Array[T]]) Array[T] {
	res := NewArray[T]()
	array.ForEach(func(v Array[T]) { res.Push(v.GetSlice()...) })
	return res
}

// ===================================快捷创建数组的方法===================================
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

func Range[T Number](start, stop T) Array[T] { return RangeStep[T](start, stop, 1) }
func RangeN[T Number](stop T) Array[T]       { return RangeStep[T](0, stop, 1) }
func RangeStep[T Number](start, stop, step T) Array[T] {
	var res []T
	for i := start; i < stop; i += step {
		res = append(res, i)
	}
	return Arr(res)
}

// 包含右边界的range函数
func RangeEq[T Number](start, stop T) Array[T] { return RangeEqStep[T](start, stop, 1) }
func RangeEqN[T Number](stop T) Array[T]       { return RangeEqStep[T](0, stop, 1) }
func RangeEqStep[T Number](start, stop, step T) Array[T] {
	var res []T
	for i := start; i <= stop; i += step {
		res = append(res, i)
	}
	return Arr(res)
}
