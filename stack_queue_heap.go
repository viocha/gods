package gods

import (
	"container/list"
	"fmt"
	"reflect"
	"slices"
)

// ===================================Heap===================================

type Heap[T any] struct {
	data   *[]T
	less   func(T, T) bool // 保证所有父元素j和子元素i满足less(j,i)
	maxCap int             // 支持TopK的堆
}

func NewHeap[T any]() Heap[T] {
	// 对可使用<=比较的类型使用自然排序，构建最小堆
	var less func(T, T) bool
	switch any(*new(T)).(type) {
	case int, int8, int16, int32, int64:
		less = func(a T, b T) bool {
			return reflect.ValueOf(a).Int() <= reflect.ValueOf(b).Int()
		}
	case uint, uint8, uint16, uint32, uint64:
		less = func(a T, b T) bool {
			return reflect.ValueOf(a).Uint() <= reflect.ValueOf(b).Uint()
		}
	case float32, float64:
		less = func(a T, b T) bool {
			return reflect.ValueOf(a).Float() <= reflect.ValueOf(b).Float()
		}
	case string:
		less = func(a T, b T) bool {
			return reflect.ValueOf(a).String() <= reflect.ValueOf(b).String()
		}
	}
	return Heap[T]{data: new([]T), less: less}
}

func (o Heap[T]) WithLess(less func(T, T) bool) Heap[T] {
	if o.Len() != 0 {
		panic("堆已有元素，不能修改less函数")
	}
	o.less = less
	return o
}

func (o Heap[T]) WithSlice(data []T) Heap[T] { // 会直接作为内部数据，不会进行拷贝
	o.data = &data
	return o.heapify()
}

func (o Heap[T]) Reversed() Heap[T] {
	less := o.less
	o.less = func(a T, b T) bool { return !less(a, b) }
	return o
}

func (o Heap[T]) WithMaxCap(maxCap int) Heap[T] {
	if o.Len() != 0 {
		panic("堆已有元素，不能修改最大容量")
	}
	o.maxCap = maxCap
	return o
}

// --------------------辅助函数--------------------

func par(i int) int { return (i - 1) / 2 }
func lc(i int) int  { return 2*i + 1 }
func rc(i int) int  { return 2*i + 2 }
func (o Heap[T]) swap(i, j int) Heap[T] {
	d := *o.data
	d[i], d[j] = d[j], d[i]
	return o
}

func (o Heap[T]) getData() (d []T, n int) { return *o.data, len(*o.data) }
func (o Heap[T]) up(i int) Heap[T] {
	if i <= 0 {
		return o
	}
	d := *o.data
	j := par(i)
	if o.less(d[j], d[i]) { // 满足父比子小
		return o
	}
	return o.swap(i, j).up(j)
}
func (o Heap[T]) down(i int) Heap[T] {
	d, n := o.getData()
	j, k := lc(i), rc(i)
	if j >= n {
		return o
	}
	if k < n && !o.less(d[j], d[k]) { // 选取最小的孩子为j
		j = k
	}
	if o.less(d[i], d[j]) { // 满足父比子小
		return o
	}
	return o.swap(i, j).down(j)
}
func (o Heap[T]) heapify() Heap[T] {
	n := len(*o.data)
	for i := par(n - 1); i >= 0; i-- {
		o.down(i)
	}
	return o
}

// --------------------其他转换--------------------

func (o Heap[T]) getSlice() []T {
	return *o.data
}

// --------------------ValueContainer接口--------------------
func (o Heap[T]) Len() int       { return len(*o.data) }
func (o Heap[T]) Clear() Heap[T] { o.data = new([]T); return o }
func (o Heap[T]) Clone() Heap[T] { return o.WithSlice(append([]T{}, *o.data...)) }
func (o Heap[T]) String() string { return "Heap" + fmt.Sprint(*o.data) }
func (o Heap[T]) ForEach(f func(T)) { // 按堆弹出顺序遍历
	h := o.Clone()
	for h.Len() > 0 {
		f(h.Pop())
	}
}
func (o Heap[T]) ToSlice() []T { // 转换成堆弹出顺序的切片
	res := make([]T, 0, o.Len())
	o.ForEach(func(x T) { res = append(res, x) })
	return res
}

// --------------------Heap接口--------------------

// 基本操作
func (o Heap[T]) Peek() T { return (*o.data)[0] }
func (o Heap[T]) Push(x T) Heap[T] {
	if o.maxCap > 0 {
		return o.PushTopK(x, o.maxCap)
	}
	return o.push(x)
}
func (o Heap[T]) Pop() T {
	d, n := o.getData()
	x := d[0]
	d[0], *o.data = d[n-1], d[:n-1]
	o.down(0)
	return x
}
func (o Heap[T]) push(x T) Heap[T] {
	d, n := o.getData()
	*o.data = append(d, x)
	return o.up(n)
}

// 组合操作
func (o Heap[T]) Set(i int, x T) Heap[T] {
	(*o.data)[i] = x
	return o.up(i).down(i)
}
func (o Heap[T]) Del(i int) Heap[T] {
	d, n := o.getData()
	d[i], *o.data = d[n-1], d[:n-1]
	return o.down(i)
}
func (o Heap[T]) PushTopK(x T, k int) Heap[T] {
	if o.Len() < k {
		return o.push(x)
	}
	if o.less(o.Peek(), x) { // x在堆顶之下
		return o.Set(0, x)
	}
	return o
}

// ===================================栈===================================
type Stack[T any] struct {
	data *[]T
}

func NewStack[T any]() Stack[T]                { return Stack[T]{new([]T)} }
func (o Stack[T]) WithSlice(data []T) Stack[T] { o.data = &data; return o }
func (o Stack[T]) getData() (d []T, n int)     { return *o.data, len(*o.data) }

// --------------------其他转换--------------------
func (o Stack[T]) GetSlice() []T { return *o.data }

// --------------------ValueContainer接口--------------------
func (o Stack[T]) Len() int        { return len(*o.data) }
func (o Stack[T]) String() string  { return "Stack" + fmt.Sprint(*o.data) }
func (o Stack[T]) Clear() Stack[T] { o.data = new([]T); return o }
func (o Stack[T]) Clone() Stack[T] { return o.WithSlice(append([]T{}, *o.data...)) }
func (o Stack[T]) ForEach(f func(T)) {
	d, n := o.getData()
	d = slices.Clone(d) // 支持边遍历边修改
	for i := n - 1; i >= 0; i-- {
		f(d[i])
	}
}
func (o Stack[T]) ToSlice() []T { return append([]T{}, *o.data...) }

// --------------------Stack接口--------------------

// 基本方法
func (o Stack[T]) Push(x T) Stack[T] { *o.data = append(*o.data, x); return o }
func (o Stack[T]) Pop() T {
	d, n := o.getData()
	x := d[n-1]
	*o.data = d[:n-1]
	return x
}
func (o Stack[T]) Peek() T {
	d, n := o.getData()
	return d[n-1]
}

// 单调栈
func (o Stack[T]) PopUntil(f func(T) bool) Stack[T] {
	for o.Len() > 0 && !f(o.Peek()) {
		o.Pop()
	}
	return o
}

// ===================================双向队列===================================

type Queue[T any] struct {
	l *list.List
}

func NewQueue[T any]() Queue[T] { return Queue[T]{list.New()} }

// --------------------ValueContainer接口--------------------
func (o Queue[T]) Len() int { return o.l.Len() }
func (o Queue[T]) Clone() Queue[T] {
	res := NewQueue[T]()
	o.ForEach(func(x T) { res.PushBack(x) })
	return res
}
func (o Queue[T]) Clear() Queue[T] { *o.l = *list.New(); return o }
func (o Queue[T]) ForEach(f func(T)) Queue[T] {
	vals := make([]T, 0)
	for e := o.l.Front(); e != nil; e = e.Next() {
		vals = append(vals, e.Value.(T))
	}
	for _, x := range vals {
		f(x)
	}
	return o
}
func (o Queue[T]) ToSlice() []T {
	res := make([]T, 0, o.Len())
	o.ForEach(func(x T) { res = append(res, x) })
	return res
}
func (o Queue[T]) String() string { return "Queue" + fmt.Sprint(o.ToSlice()) }

// --------------------Queue接口--------------------

// 单向队列操作
func (o Queue[T]) Push(x T) Queue[T] { return o.PushBack(x) }
func (o Queue[T]) Pop() T            { return o.PopFront() }
func (o Queue[T]) Peek() T           { return o.PeekFront() }

// 双向队列操作
func (o Queue[T]) PushFront(x T) Queue[T] { o.l.PushFront(x); return o }
func (o Queue[T]) PushBack(x T) Queue[T]  { o.l.PushBack(x); return o }
func (o Queue[T]) PopFront() T            { return o.l.Remove(o.l.Front()).(T) }
func (o Queue[T]) PopBack() T             { return o.l.Remove(o.l.Back()).(T) }
func (o Queue[T]) PeekFront() T           { return o.l.Front().Value.(T) }
func (o Queue[T]) PeekBack() T            { return o.l.Back().Value.(T) }

// 单调队列
func (o Queue[T]) PopUntil(f func(T) bool) Queue[T] { return o.PopFrontUntil(f) }
func (o Queue[T]) PopFrontUntil(f func(T) bool) Queue[T] {
	for o.Len() > 0 && !f(o.PeekFront()) {
		o.PopFront()
	}
	return o
}
func (o Queue[T]) PopBackUntil(f func(T) bool) Queue[T] {
	for o.Len() > 0 && !f(o.PeekBack()) {
		o.PopBack()
	}
	return o
}
