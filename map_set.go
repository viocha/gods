package main

import (
	"fmt"
	"maps"
	"strings"
)

// ===================================哈希表===================================

type HashMap[K comparable, V any] struct {
	m       map[K]V
	factory func() V
}

func NewHashMap[K comparable, V any]() HashMap[K, V] { return HashMap[K, V]{m: map[K]V{}} }
func (o HashMap[K, V]) WithFactory(factory func() V) HashMap[K, V] {
	o.factory = factory
	return o
}
func (o HashMap[K, V]) WithMap(m map[K]V) HashMap[K, V] { o.m = m; return o }

// --------------------其他转换--------------------

func (o HashMap[K, V]) GetMap() map[K]V { return o.m }

// --------------------Container接口--------------------

func (o HashMap[K, V]) Len() int             { return len(o.m) }
func (o HashMap[K, V]) Clear() HashMap[K, V] { clear(o.m); return o }
func (o HashMap[K, V]) Clone() HashMap[K, V] { return o.WithMap(maps.Clone(o.m)) }
func (o HashMap[K, V]) String() string       { return strings.ReplaceAll(fmt.Sprint(o.m), "map", "HashMap") }

// --------------------Map接口--------------------
// 遍历和转换
func (o HashMap[K, V]) ForEach(f func(K, V)) HashMap[K, V] {
	for k, v := range o.m {
		f(k, v)
	}
	return o
}
func (o HashMap[K, V]) ToMap() map[K]V { return maps.Clone(o.m) }

// 基本操作
func (o HashMap[K, V]) Has(k K) bool               { _, ok := o.m[k]; return ok }
func (o HashMap[K, V]) Set(k K, v V) HashMap[K, V] { o.m[k] = v; return o }
func (o HashMap[K, V]) Del(k K) HashMap[K, V]      { delete(o.m, k); return o }
func (o HashMap[K, V]) Get(k K) V {
	if !o.Has(k) && o.factory != nil {
		o.Set(k, o.factory())
	}
	return o.m[k]
}

// 组合操作
func (o HashMap[K, V]) GetOr(k K, v V) V {
	if !o.Has(k) {
		return v
	}
	return o.m[k]
}
func (o HashMap[K, V]) GetOrSet(k K, v V) V {
	if !o.Has(k) {
		o.Set(k, v)
	}
	return o.m[k]
}
func (o HashMap[K, V]) Extend(other HashMap[K, V]) HashMap[K, V] {
	other.ForEach(func(k K, v V) { o.Set(k, v) })
	return o
}
func (o HashMap[K, V]) DelFunc(f func(K, V) bool) HashMap[K, V] {
	return o.ForEach(func(k K, v V) {
		if f(k, v) {
			o.Del(k)
		}
	})
}
func (o HashMap[K, V]) ReplaceFunc(f func(K, V) V) HashMap[K, V] {
	return o.ForEach(func(k K, v V) { o.Set(k, f(k, v)) })
}

// 键值查询
func (o HashMap[K, V]) Keys() []K {
	var keys []K
	o.ForEach(func(k K, v V) { keys = append(keys, k) })
	return keys
}
func (o HashMap[K, V]) Values() []V {
	var values []V
	o.ForEach(func(k K, v V) { values = append(values, v) })
	return values
}

// ===================================哈希集===================================
type HashSet[T comparable] struct {
	m HashMap[T, struct{}]
}

func NewHashSet[T comparable]() HashSet[T] { return HashSet[T]{m: NewHashMap[T, struct{}]()} }
func NewHashSetFromSlice[T comparable](s []T) HashSet[T] {
	res := NewHashSet[T]()
	for _, v := range s {
		res.Add(v)
	}
	return res
}

// --------------------ValueContainer接口--------------------
func (o HashSet[T]) Len() int          { return o.m.Len() }
func (o HashSet[T]) Clear() HashSet[T] { o.m.Clear(); return o }
func (o HashSet[T]) Clone() HashSet[T] { o.m = o.m.Clone(); return o }
func (o HashSet[T]) String() string    { return "HashSet" + fmt.Sprint(o.m.Keys()) }
func (o HashSet[T]) ToSlice() []T      { return o.m.Keys() }
func (o HashSet[T]) ForEach(f func(T)) HashSet[T] {
	o.m.ForEach(func(k T, v struct{}) { f(k) })
	return o
}

// --------------------Set接口--------------------

// 基本操作
func (o HashSet[T]) Has(v T) bool       { return o.m.Has(v) }
func (o HashSet[T]) Del(v T) HashSet[T] { o.m.Del(v); return o }
func (o HashSet[T]) Add(v T) HashSet[T] { o.m.Set(v, struct{}{}); return o }

// 组合操作
func (o HashSet[T]) HasSubset(other HashSet[T]) bool {
	for _, v := range other.ToSlice() {
		if !o.Has(v) {
			return false
		}
	}
	return true
}
func (o HashSet[T]) DelFunc(f func(T) bool) HashSet[T] {
	return o.ForEach(func(v T) {
		if f(v) {
			o.Del(v)
		}
	})
}
func (o HashSet[T]) Union(other HashSet[T]) HashSet[T] {
	res := o.Clone()
	other.ForEach(func(v T) { res.Add(v) })
	return res
}
func (o HashSet[T]) Intersect(other HashSet[T]) HashSet[T] {
	res := NewHashSet[T]()
	o.ForEach(func(v T) {
		if other.Has(v) {
			res.Add(v)
		}
	})
	return res
}
func (o HashSet[T]) Difference(other HashSet[T]) HashSet[T] {
	res := NewHashSet[T]()
	o.ForEach(func(v T) {
		if !other.Has(v) {
			res.Add(v)
		}
	})
	return res
}

// ===================================链式哈希表===================================

// --------------------链式哈希表所需的链表结构--------------------

type linkedListNode[K comparable, V any] struct {
	k          K
	v          V
	prev, next *linkedListNode[K, V]
}
type linkedList[K comparable, V any] struct {
	root *linkedListNode[K, V]
}

func newLinkedList[K comparable, V any]() linkedList[K, V] {
	root := &linkedListNode[K, V]{}
	root.prev, root.next = root, root
	return linkedList[K, V]{root: root}
}
func (o linkedList[K, V]) clear() linkedList[K, V] {
	root := o.root
	root.next, root.prev = root, root
	return o
}
func (o linkedList[K, V]) front() *linkedListNode[K, V] {
	if o.root.next == o.root {
		return nil
	}
	return o.root.next
}

func (o linkedList[K, V]) addBack(node *linkedListNode[K, V]) linkedList[K, V] {
	pre := o.root.prev
	node.prev, node.next = pre, o.root
	pre.next, o.root.prev = node, node
	return o
}
func (o linkedList[K, V]) remove(node *linkedListNode[K, V]) linkedList[K, V] {
	pre, nex := node.prev, node.next
	node.prev, node.next = nil, nil
	pre.next, nex.prev = nex, pre
	return o
}
func (o linkedList[K, V]) moveToBack(node *linkedListNode[K, V]) linkedList[K, V] {
	return o.remove(node).addBack(node)
}

// --------------------链式哈希表--------------------

type LinkedHashMap[K comparable, V any] struct {
	m           map[K]*linkedListNode[K, V]
	l           linkedList[K, V]
	accessOrder bool
	factory     func() V
	maxCap      int
}

func NewLinkedHashMap[K comparable, V any]() LinkedHashMap[K, V] {
	return LinkedHashMap[K, V]{m: make(map[K]*linkedListNode[K, V]), l: newLinkedList[K, V]()}
}

func (o LinkedHashMap[K, V]) WithFactory(factory func() V) LinkedHashMap[K, V] {
	o.factory = factory
	return o
}

func (o LinkedHashMap[K, V]) WithAccessOrderMode() LinkedHashMap[K, V] {
	if o.Len() != 0 {
		panic("哈希表不为空，不能设置accessOrder模式")
	}
	o.accessOrder = true
	return o
}

func (o LinkedHashMap[K, V]) WithMaxCap(maxCap int) LinkedHashMap[K, V] {
	if o.Len() != 0 {
		panic("哈希表不为空，不能设置最大容量")
	}
	o.maxCap = maxCap
	return o
}

// --------------------Container接口--------------------
func (o LinkedHashMap[K, V]) Len() int { return len(o.m) }
func (o LinkedHashMap[K, V]) Clear() LinkedHashMap[K, V] {
	clear(o.m)
	o.l.clear()
	return o
}
func (o LinkedHashMap[K, V]) Clone() LinkedHashMap[K, V] {
	o.m = make(map[K]*linkedListNode[K, V])
	o.l = newLinkedList[K, V]()
	o.ForEach(func(k K, v V) { o.Set(k, v) })
	return o
}
func (o LinkedHashMap[K, V]) String() string {
	entries := make([]string, 0)
	o.ForEach(func(k K, v V) { entries = append(entries, fmt.Sprintf("%v:%v", k, v)) })
	return "LinkedHashMap" + fmt.Sprint(entries)
}

// --------------------Map接口--------------------
// 遍历和转换
func (o LinkedHashMap[K, V]) ForEach(f func(K, V)) LinkedHashMap[K, V] {
	if o.Len() == 0 {
		return o
	}
	nodes := make([]linkedListNode[K, V], 0, o.Len())
	for p := o.l.front(); p != o.l.root; p = p.next {
		nodes = append(nodes, *p)
	}
	for _, node := range nodes {
		f(node.k, node.v)
	}
	return o
}
func (o LinkedHashMap[K, V]) ToMap() map[K]V {
	res := make(map[K]V, o.Len())
	o.ForEach(func(k K, v V) { res[k] = v })
	return res
}

// 基本操作
func (o LinkedHashMap[K, V]) Has(k K) bool { _, ok := o.m[k]; return ok }
func (o LinkedHashMap[K, V]) Set(k K, v V) LinkedHashMap[K, V] {
	if !o.Has(k) { // 添加
		o.ensureCap()
		node := &linkedListNode[K, V]{k: k, v: v}
		o.l.addBack(node)
		o.m[k] = node
		return o
	}
	node := o.m[k] // 修改
	node.v = v
	o.afterAccess(node)
	return o
}
func (o LinkedHashMap[K, V]) Get(k K) V {
	if !o.Has(k) && o.factory != nil {
		o.Set(k, o.factory())
	}
	node := o.m[k]
	o.afterAccess(node)
	return node.v
}
func (o LinkedHashMap[K, V]) Del(k K) LinkedHashMap[K, V] {
	node := o.m[k]
	delete(o.m, k)
	o.l.remove(node)
	return o
}

func (o LinkedHashMap[K, V]) afterAccess(node *linkedListNode[K, V]) LinkedHashMap[K, V] {
	if o.accessOrder {
		o.l.moveToBack(node)
	}
	return o
}
func (o LinkedHashMap[K, V]) ensureCap() LinkedHashMap[K, V] {
	if o.maxCap <= 0 {
		return o
	}
	if o.Len() == o.maxCap {
		o.Del(o.l.front().k)
	}
	return o
}

// 组合操作
func (o LinkedHashMap[K, V]) GetOrSet(k K, v V) V {
	if !o.Has(k) {
		o.Set(k, v)
	}
	return o.Get(k)
}
func (o LinkedHashMap[K, V]) GetOr(k K, v V) V {
	if !o.Has(k) {
		return v
	}
	return o.Get(k)
}
func (o LinkedHashMap[K, V]) Extend(other LinkedHashMap[K, V]) LinkedHashMap[K, V] {
	other.ForEach(func(k K, v V) { o.Set(k, v) })
	return o
}
func (o LinkedHashMap[K, V]) DelFunc(f func(K, V) bool) LinkedHashMap[K, V] {
	return o.ForEach(func(k K, v V) {
		if f(k, v) {
			o.Del(k)
		}
	})
}
func (o LinkedHashMap[K, V]) ReplaceFunc(f func(K, V) V) LinkedHashMap[K, V] {
	return o.ForEach(func(k K, v V) { o.Set(k, f(k, v)) })
}

// 键值查询
func (o LinkedHashMap[K, V]) Keys() []K {
	res := make([]K, 0, o.Len())
	o.ForEach(func(k K, v V) { res = append(res, k) })
	return res
}
func (o LinkedHashMap[K, V]) Values() []V {
	res := make([]V, 0, o.Len())
	o.ForEach(func(k K, v V) { res = append(res, v) })
	return res
}

// --------------------LinkedMapContainer接口--------------------
func (o LinkedHashMap[K, V]) FirstKey() K {
	if o.Len() == 0 {
		panic("哈希表为空")
	}
	return o.l.front().k
}
func (o LinkedHashMap[K, V]) LastKey() K {
	if o.Len() == 0 {
		panic("哈希表为空")
	}
	return o.l.root.prev.k
}
func (o LinkedHashMap[K, V]) NextKey(k K) K {
	if o.Len() == 0 {
		panic("哈希表为空")
	}
	node := o.m[k]
	if node.next == o.l.root {
		panic("最后一个节点，没有下一个节点")
	}
	return node.next.k
}
func (o LinkedHashMap[K, V]) PrevKey(k K) K {
	if o.Len() == 0 {
		panic("哈希表为空")
	}
	node := o.m[k]
	if node.prev == o.l.root {
		panic("第一个节点，没有前一个节点")
	}
	return node.prev.k
}

// ===================================链式哈希集===================================

type LinkedHashSet[T comparable] struct {
	m LinkedHashMap[T, struct{}]
}

func NewLinkedHashSet[T comparable]() LinkedHashSet[T] {
	return LinkedHashSet[T]{m: NewLinkedHashMap[T, struct{}]()}
}

func (o LinkedHashSet[T]) WithAccessOrderMode() LinkedHashSet[T] {
	o.m = o.m.WithAccessOrderMode()
	return o
}
func (o LinkedHashSet[T]) WithMaxCap(maxCap int) LinkedHashSet[T] {
	o.m = o.m.WithMaxCap(maxCap)
	return o
}

func NewLinkedHashSetFromSlice[T comparable](s []T) LinkedHashSet[T] {
	o := NewLinkedHashSet[T]()
	for _, v := range s {
		o.Add(v)
	}
	return o
}

// --------------------ValueContainer接口--------------------

func (o LinkedHashSet[T]) Len() int                { return o.m.Len() }
func (o LinkedHashSet[T]) Clear() LinkedHashSet[T] { o.m.Clear(); return o }
func (o LinkedHashSet[T]) Clone() LinkedHashSet[T] { o.m = o.m.Clone(); return o }
func (o LinkedHashSet[T]) ForEach(f func(T)) LinkedHashSet[T] {
	o.m.ForEach(func(k T, v struct{}) { f(k) })
	return o
}
func (o LinkedHashSet[T]) ToSlice() []T { return o.m.Keys() }
func (o LinkedHashSet[T]) String() string {
	return "LinkedHashSet" + fmt.Sprint(o.ToSlice())
}

// --------------------SetContainer接口--------------------

// 基本操作
func (o LinkedHashSet[T]) Has(k T) bool { // 查询是否存在会将元素移动到链表尾部
	has := o.m.Has(k)
	if has {
		o.m.Get(k)
	}
	return has
}
func (o LinkedHashSet[T]) Add(k T) LinkedHashSet[T] {
	o.m.Set(k, struct{}{})
	return o
}
func (o LinkedHashSet[T]) Del(k T) LinkedHashSet[T] {
	o.m.Del(k)
	return o
}

// 组合操作
func (o LinkedHashSet[T]) HasSubset(other LinkedHashSet[T]) bool {
	for _, v := range other.ToSlice() {
		if !o.Has(v) {
			return false
		}
	}
	return true
}
func (o LinkedHashSet[T]) DelFunc(f func(T) bool) LinkedHashSet[T] {
	return o.ForEach(func(k T) {
		if f(k) {
			o.Del(k)
		}
	})
}
func (o LinkedHashSet[T]) Union(other LinkedHashSet[T]) LinkedHashSet[T] {
	res := o.Clone()
	other.ForEach(func(k T) { res.Add(k) })
	return res
}
func (o LinkedHashSet[T]) Intersect(other LinkedHashSet[T]) LinkedHashSet[T] {
	res := NewLinkedHashSet[T]()
	o.ForEach(func(x T) {
		if other.Has(x) {
			res.Add(x)
		}
	})
	return res
}
func (o LinkedHashSet[T]) Difference(other LinkedHashSet[T]) LinkedHashSet[T] {
	res := NewLinkedHashSet[T]()
	o.ForEach(func(x T) {
		if !other.Has(x) {
			res.Add(x)
		}
	})
	return res
}

// --------------------LinkedSetContainer接口--------------------

func (o LinkedHashSet[T]) First() T   { return o.m.FirstKey() }
func (o LinkedHashSet[T]) Last() T    { return o.m.LastKey() }
func (o LinkedHashSet[T]) Next(k T) T { return o.m.NextKey(k) }
func (o LinkedHashSet[T]) Prev(k T) T { return o.m.PrevKey(k) }

// ===================================多重哈希集===================================
type MultiHashSet[T comparable] struct {
	m     HashMap[T, int]
	total *int // 重数总和，仅保存引用，以便值传递进行修改
}

func NewMultiHashSet[T comparable]() MultiHashSet[T] {
	return MultiHashSet[T]{
		m:     NewHashMap[T, int]().WithFactory(func() int { return 0 }),
		total: new(int),
	}
}
func NewMultiHashSetFromSlice[T comparable](s []T) MultiHashSet[T] {
	o := NewMultiHashSet[T]()
	for _, v := range s {
		o.Add(v)
	}
	return o
}
func NewMultiHashSetFromMap[T comparable](m map[T]int) MultiHashSet[T] {
	o := NewMultiHashSet[T]()
	for k, v := range m {
		o.m.Set(k, v)
	}
	return o
}

// --------------------其他转换--------------------
func (o MultiHashSet[T]) ToHashSet() HashSet[T] { return NewHashSetFromSlice(o.ToSlice()) } // 去重

// --------------------ValueContainer接口--------------------
func (o MultiHashSet[T]) Len() int { return o.m.Len() }
func (o MultiHashSet[T]) ForEach(f func(T)) MultiHashSet[T] {
	o.ForEachCnt(func(x T, cnt int) { f(x) })
	return o
}
func (o MultiHashSet[T]) ToSlice() []T {
	res := make([]T, 0, o.Total())
	o.ForEachCnt(func(x T, cnt int) {
		for range cnt {
			res = append(res, x)
		}
	})
	return res
}
func (o MultiHashSet[T]) String() string {
	return strings.ReplaceAll(o.m.String(), "HashMap", "MultiHashSet")
}
func (o MultiHashSet[T]) Clear() MultiHashSet[T] { o.m.Clear(); return o }
func (o MultiHashSet[T]) Clone() MultiHashSet[T] { return NewMultiHashSetFromSlice(o.ToSlice()) }

// --------------------SetContainer接口--------------------

// 基本操作
func (o MultiHashSet[T]) Has(x T) bool            { return o.m.Has(x) }
func (o MultiHashSet[T]) Add(x T) MultiHashSet[T] { return o.AddN(x, 1) }
func (o MultiHashSet[T]) Del(x T) MultiHashSet[T] { return o.DelN(x, 1) }

// 组合操作
func (o MultiHashSet[T]) HasSubset(other MultiHashSet[T]) bool {
	for v, cnt := range other.m.GetMap() {
		if cnt > o.Count(v) {
			return false
		}
	}
	return true
}
func (o MultiHashSet[T]) DelFunc(f func(T) bool) MultiHashSet[T] {
	return o.ForEach(func(v T) {
		if f(v) {
			o.Del(v)
		}
	})
}
func (o MultiHashSet[T]) Intersect(other MultiHashSet[T]) MultiHashSet[T] {
	res := NewMultiHashSet[T]()
	other.ForEachCnt(func(v T, cnt int) { res.AddN(v, min(cnt, o.Count(v))) })
	return res
}
func (o MultiHashSet[T]) Union(other MultiHashSet[T]) MultiHashSet[T] {
	res := o.Clone()
	other.ForEachCnt(func(v T, cnt int) { res.AddN(v, cnt-o.Count(v)) })
	return res
}
func (o MultiHashSet[T]) Difference(other MultiHashSet[T]) MultiHashSet[T] {
	res := o.Clone()
	other.ForEachCnt(func(v T, cnt int) { res.DelN(v, o.Count(v)-cnt) })
	return res
}

// --------------------MultiSetContainer接口--------------------

// 遍历和转换
func (o MultiHashSet[T]) ForEachCnt(f func(T, int)) MultiHashSet[T] { o.m.ForEach(f); return o }
func (o MultiHashSet[T]) ToMap() map[T]int {
	res := make(map[T]int, o.Len())
	o.ForEachCnt(func(x T, cnt int) { res[x] = cnt })
	return res
}

// 计数操作
func (o MultiHashSet[T]) AddN(x T, n int) MultiHashSet[T] {
	if n <= 0 {
		return o
	}
	o.m.Set(x, o.m.Get(x)+n)
	*o.total += n
	return o
}
func (o MultiHashSet[T]) DelN(x T, n int) MultiHashSet[T] {
	if n <= 0 || !o.Has(x) {
		return o
	}
	cnt := o.m.Get(x)
	if n >= cnt {
		o.m.Del(x)
		*o.total -= cnt
	} else {
		o.m.Set(x, cnt-n)
		*o.total -= n
	}
	return o
}
func (o MultiHashSet[T]) DelAll(x T) MultiHashSet[T] { return o.DelN(x, o.Count(x)) }
func (o MultiHashSet[T]) Total() int                 { return *o.total }
func (o MultiHashSet[T]) Count(x T) int {
	if !o.Has(x) {
		return 0
	}
	return o.m.Get(x)
}
func (o MultiHashSet[T]) DelAllFunc(f func(T, int) bool) MultiHashSet[T] {
	return o.ForEachCnt(func(x T, cnt int) {
		if f(x, cnt) {
			o.DelAll(x)
		}
	})
}
func (o MultiHashSet[T]) ReplaceFunc(f func(T, int) int) MultiHashSet[T] {
	return o.ForEachCnt(func(x T, cnt int) {
		newCnt := f(x, cnt)
		if newCnt > cnt {
			o.AddN(x, newCnt-cnt)
		} else if newCnt < cnt {
			o.DelN(x, cnt-newCnt)
		}
	})
}
