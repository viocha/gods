package main

import "cmp"

/* ===========================通用接口============================ */
type Container[T any] interface {
	Len() int
	String() string
	Clear() Container[T]
	Clone() Container[T]
}

type TreePrinter interface {
	PrintTree()
}
type Slicer[T any] interface {
	ToSlice() []T
}
type Mapper[K comparable, V any] interface {
	ToMap() map[K]V
}

/* ===========================单值容器============================ */
// 可以转换成slice的容器
type ValueContainer[T any] interface {
	Container[T]
	Slicer[T]
	ForEach(f func(T)) ValueContainer[T]
}

// 堆
type HeapContainer[T any] interface { // 支持自然排序和自定义排序
	ValueContainer[T]
	Reversed() HeapContainer[T] // 将堆的大小关系取反，并重新建堆
	// 基本操作
	Pop() T
	Peek() T
	Push(v T) HeapContainer[T]
	// 组合操作
	Set(i int, v T) HeapContainer[T]
	Del(i int, v T) HeapContainer[T]
	PushTopK(v T, k int) HeapContainer[T] // 维持堆大小为k，获取最小k个元素
}

// 栈
type StackContainer[T any] interface {
	ValueContainer[T]
	// 基本操作
	Push(v T) StackContainer[T]
	Pop() T
	Peek() T
	// 单调栈操作
	PopUntil(f func(T) bool) StackContainer[T]
}

// 队列
type QueueContainer[T any] interface {
	ValueContainer[T]
	// 单向队列基本操作
	Push(v T) QueueContainer[T]
	Pop() T
	Peek() T
	// 双向队列基本操作
	Front() T
	Back() T
	PushFront(v T) QueueContainer[T]
	PushBack(v T) QueueContainer[T]
	PopFront() T
	PopBack() T
	// 单调队列操作
	PopUntil(f func(T) bool) QueueContainer[T]
	PopFrontUntil(f func(T) bool) QueueContainer[T]
	PopBackUntil(f func(T) bool) QueueContainer[T]
}

// 环形数组 TODO
type RingContainer[T any] interface {
}

// ==============基本容器=============

// 动态数组
type ArrayContainer[T any] interface {
	ValueContainer[T]
	// 遍历
	ForEachIdx(f func(int)) ArrayContainer[T]       // 遍历索引
	ForEachIdxVal(f func(int, T)) ArrayContainer[T] // 遍历索引和值

	// 基本操作(支持负数索引)
	Get(i int) T
	Set(i int, v T) ArrayContainer[T]
	Ins(i int, v T) ArrayContainer[T]
	Del(i int) ArrayContainer[T]
	Push(v ...T) ArrayContainer[T]     // 从尾部插入
	Pop() T                            // 从尾部删除
	Slice(i, j int) ArrayContainer[T]  // 修改成切片内容
	Sliced(i, j int) ArrayContainer[T] // 返回新的切片视图，不修改内部数组

	// 组合操作
	Replace(i, j int, v ...T) ArrayContainer[T]     // 将i到j之间的元素替换成1个或多个元素
	ReplaceFunc(f func(int, T) T) ArrayContainer[T] // 替换所有元素的值
	DelFunc(f func(T) bool) ArrayContainer[T]       // 删除满足条件的所有元素
	Extend(v ArrayContainer[T]) ArrayContainer[T]   // 在尾部添加切片中所有元素
	Fill(v T, n int) ArrayContainer[T]
	Reverse() ArrayContainer[T]  // 原地反转
	Reversed() ArrayContainer[T] // 返回反转的新数组

	// 查找、排序、比较
	Has(x T) bool // 查找指定元素位置
	HasBi(x T) bool
	Index(v T) int
	IndexFunc(eq func(T) bool) int // 线性查找满足条件的索引
	IndexBi(x T) int               // 当T是comparable时，线性查找
	IndexBiFunc(x T, cmp func(T, T) int) int

	HasSubArr(v ArrayContainer[T]) bool // 查找子数组位置
	HasSubArrFunc(v ArrayContainer[T], eq func(T, T) bool) bool
	IndexSubArr(v ArrayContainer[T]) int
	IndexSubArrFunc(v ArrayContainer[T], eq func(T, T) bool) int

	BiSearch(v T) int                    // 二分查找>=v的位置
	BiSearchFunc(check func(T) bool) int // 二分查找第一个满足条件的位置

	Sort() ArrayContainer[T]                       // 当T是cmp.Ordered时，原地排序
	SortFunc(cmp func(T, T) int) ArrayContainer[T] // 自定义排序，默认稳定排序

	Cmp(other ArrayContainer[T]) int // 按元素偏序比较
	CmpFunc(f func(T, T) int) int
	Equal(other ArrayContainer[T]) bool // 按元素相等比较
	EqualFunc(f func(T, T) bool) bool

	// 随机化操作
	Shuffle() ArrayContainer[T] // 原地打乱
	Choice() ArrayContainer[T]  // 随机选择一个元素

	// 聚合操作
	Min() T                        // 自然排序的最小值
	MinFunc(less func(T, T) int) T // 自定义排序的最小值
	Max() T                        // 自然排序的最大值
	MaxFunc(less func(T, T) int) T // 自定义排序的最大值
	Sum() T                        // 对可以相加的元素求和

	Any() bool // 和零值判断，可以先Map成布尔
	All() bool
	Join(seq string) string // 使用Sprintf转换成字符串，然后拼接

	// 函数式操作(非原地)
	Map(f func(T) T) ArrayContainer[T]
	MapTo(f func(T) any) ArrayContainer[any]
	MapToInt(f func(T) int) ArrayContainer[int]
	MapToFloat(f func(T) float64) ArrayContainer[float64]
	MapToStr(f func(T) string) ArrayContainer[string]
	MapToBool(f func(T) bool) ArrayContainer[bool]
	Filter(f func(T) bool) ArrayContainer[T]
	Reduce(f func(T, T) T) T
	ReduceInitial(f func(T, T) T, initial T) T
}

// 双向链表
type ListContainer[T any] interface {
	ValueContainer[T]

	// 基本操作
	Front() any
	Back() any
	PushFront(v T) any // 返回新建的节点
	PushBack(v T) any
	InsAfter(v T, mark any) any
	InstBefore(v T, mark any) any
	PopFront() any
	PopBack() any
	Remove(node any) T // 删除后返回节点值

	// 移动节点
	MoveBefore(node, mark any) T
	MoveAfter(node, mark any) T
	MoveToFront(node any) T
	MoveToBack(node any) T
}

// 字符串包装 TODO
type StringContainer interface {
}

// --------------------集合--------------------
type SetContainer[T comparable] interface {
	ValueContainer[T]

	// 集合基本操作
	Add(v T) SetContainer[T]
	Del(v T) SetContainer[T]
	Has(v T) bool

	// 组合操作
	DelFunc(f func(T) bool) SetContainer[T]
	HasSubset(s SetContainer[T]) bool
	Union(s SetContainer[T]) SetContainer[T]
	Intersect(s SetContainer[T]) SetContainer[T]
	Difference(s SetContainer[T]) SetContainer[T]
}

// 有序集合
type TreeSetContainer[T cmp.Ordered] interface {
	SetContainer[T]

	// 二分查找值，可能panic，需要确保存在
	First() T
	Last() T
	Lower(x T) T
	Upper(x T) T
	Floor(x T) T
	Ceiling(x T) T

	// 截取子集合
	HeadSet(k int) TreeSetContainer[T]
	TailSet(k int) TreeSetContainer[T]

	// 排名相关的查找
	Rank(x T) int
	Select(k int) T
}

// 多重集
type MultiSetContainer[T comparable] interface {
	SetContainer[T]
	Mapper[T, int] // ToMap方法

	// 计数操作
	ForEachCnt(f func(T, int))            // 遍历元素及其重数
	Total() int                           // 所有重数之和
	Count(x T) int                        // 元素的重数
	AddN(x T, n int) MultiSetContainer[T] // AddN忽略小于等于0的数
	DelN(x T, n int) MultiSetContainer[T] // DelN忽略小于等于0的数
	DelAll(x T)                           // 删除元素的所有重数
	DelAllFunc(f func(T, int) bool)       // 删除满足条件的元素的所有重数
	ReplaceFunc(f func(T, int) int)       // 替换元素的重数
}

// 链式集合
type LinkedSetContainer[T comparable] interface {
	SetContainer[T] // 不支持交并补操作

	// 链式操作
	First() T
	Last() T
	Prev(v T) T
	Next(v T) T
}

/* ===========================键值对容器============================ */

// 通用映射
type MapContainer[K comparable, V any] interface {
	// 转换和遍历
	Container[K]
	Mapper[K, V]
	ForEach(f func(K, V)) MapContainer[K, V]

	// 基本操作
	Get(k K) V // 获取键对应的值，对于不存在的key，若factory属性不为空，则使用factory函数设置值，否则返回类型零值(不会报错)
	Has(k K) bool
	Set(k K, v V) MapContainer[K, V]
	Del(k K) MapContainer[K, V]

	// 组合操作
	GetOr(k K, v V) V
	GetOrSet(k K, v V) V
	Extend(m MapContainer[K, V]) MapContainer[K, V]
	DelFunc(f func(K, V) bool) MapContainer[K, V]
	ReplaceFunc(f func(K, V) V) MapContainer[K, V] // 替换键值对的值

	// 键值列表查询
	Keys() []K
	Values() []V
}

// 有序映射
type TreeMapContainer[K cmp.Ordered, V any] interface {
	MapContainer[K, V]
	// 二分查找，可能返回nil
	First() any
	Last() any
	Lower(x K) any
	Upper(x K) any
	Floor(x K) any
	Ceiling(x K) any

	// 截取子映射
	HeadMap(k K) TreeMapContainer[K, V]
	TailMap(k K) TreeMapContainer[K, V]

	// 排名相关的查找
	Rank(x K) int
	Select(k int) any
}

// 能保存插入顺序的映射
type LinkedMapContainer[K comparable, V any] interface {
	MapContainer[K, V]

	// 链式操作
	FirstKey() K
	LastKey() K
	NextKey(k K) K
	PrevKey(k K) K
}
