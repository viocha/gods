package gods

import (
	"cmp"
	"fmt"
	"maps"
	"strings"
)

// ===================================键值对记录===================================

type MapEntry[K cmp.Ordered, V any] struct {
	node *avlNode[K, V]
	K    K
	V    V
}

func newEntryFromNode[K cmp.Ordered, V any](node *avlNode[K, V]) *MapEntry[K, V] {
	if node == nil {
		return nil
	}
	return &MapEntry[K, V]{
		node: node,
		K:    node.k,
		V:    node.v,
	}
}

// 左右节点对应的Entry
func (o *MapEntry[K, V]) Left() *MapEntry[K, V]  { return newEntryFromNode(o.node.left) }
func (o *MapEntry[K, V]) Right() *MapEntry[K, V] { return newEntryFromNode(o.node.right) }

// ===================================监视节点变化的函数===================================
type Watcher[K cmp.Ordered, V any] func(cur, left, right *MapEntry[K, V]) // 用于额外记录和统计节点的信息，比如求和

// ===================================AVL树节点===================================

type avlNode[K cmp.Ordered, V any] struct {
	k           K
	v           V
	left, right *avlNode[K, V]
	depth       int // 子树高度
	size        int // 子树节点个数
}

func newAVLNode[K cmp.Ordered, V any](k K, v V) *avlNode[K, V] {
	return &avlNode[K, V]{
		k:     k,
		v:     v,
		depth: 1,
		size:  1,
	}
}

func (o avlNode[K, V]) factor() int {
	ld, rd := 0, 0
	if o.left != nil {
		ld = o.left.depth
	}
	if o.right != nil {
		rd = o.right.depth
	}
	return rd - ld
}
func (o *avlNode[K, V]) updateStatus(w Watcher[K, V]) *avlNode[K, V] {
	ld, rd := 0, 0
	if o.left != nil {
		ld = o.left.depth
	}
	if o.right != nil {
		rd = o.right.depth
	}
	lSize, rSize := 0, 0
	if o.left != nil {
		lSize = o.left.size
	}
	if o.right != nil {
		rSize = o.right.size
	}
	o.depth = max(ld, rd) + 1
	o.size = lSize + rSize + 1
	if w != nil {
		w(newEntryFromNode(o), newEntryFromNode(o.left), newEntryFromNode(o.right))
	}
	return o
}
func (o *avlNode[K, V]) rotateLeft(w Watcher[K, V]) *avlNode[K, V] {
	r := o.right
	rl := r.left

	r.left = o
	o.right = rl
	o.updateStatus(w)
	r.updateStatus(w)
	return r
}
func (o *avlNode[K, V]) rotateRight(w Watcher[K, V]) *avlNode[K, V] {
	l := o.left
	lr := l.right

	l.right = o
	o.left = lr
	o.updateStatus(w)
	l.updateStatus(w)
	return l
}
func (o *avlNode[K, V]) balance(w Watcher[K, V]) *avlNode[K, V] {
	if o.factor() == -2 {
		if o.left.factor() > 0 { // LR
			o.left = o.left.rotateLeft(w)
		}
		return o.rotateRight(w) // LL
	} else if o.factor() == 2 {
		if o.right.factor() < 0 { // RL
			o.right = o.right.rotateRight(w)
		}
		return o.rotateLeft(w) // RR
	}
	return o
}

// ===================================基于AVL树的有序映射===================================
type TreeMap[K cmp.Ordered, V any] struct {
	dummyRoot *avlNode[K, V] // 实际的根节点为dummyRoot.right
	factory   func() V
	w         Watcher[K, V]
}

func NewTreeMap[K cmp.Ordered, V any]() TreeMap[K, V] {
	return TreeMap[K, V]{dummyRoot: newAVLNode(*new(K), *new(V))}
}
func NewTreeMapFromMap[K cmp.Ordered, V any](m map[K]V) TreeMap[K, V] {
	o := NewTreeMap[K, V]()
	for k, v := range m {
		o.Set(k, v)
	}
	return o
}

// 创建具有默认值的映射，在调用get时，若key不存在，则使用factory函数设置值
func (o TreeMap[K, V]) WithFactory(factory func() V) TreeMap[K, V] {
	o.factory = factory
	return o
}

// 给TreeMap添加Watcher，必须在映射的长度为0且未设置Watcher时调用
func (o TreeMap[K, V]) WithWatcher(w Watcher[K, V]) TreeMap[K, V] {
	if o.Len() != 0 {
		panic("映射必须空，才能设置Watcher")
	}
	o.w = w
	return o
}

func (o TreeMap[K, V]) GetRoot() *MapEntry[K, V] {
	return newEntryFromNode(o.dummyRoot.right)
}

func (o TreeMap[K, V]) root() *avlNode[K, V]                      { return o.dummyRoot.right }
func (o TreeMap[K, V]) setRoot(root *avlNode[K, V]) TreeMap[K, V] { o.dummyRoot.right = root; return o }

// ==============Container接口============= */
func (o TreeMap[K, V]) Len() int {
	if o.root() == nil {
		return 0
	}
	return o.root().size
}
func (o TreeMap[K, V]) String() string       { return "TreeMap" + fmt.Sprint(o.ToMap())[3:] }
func (o TreeMap[K, V]) Clear() TreeMap[K, V] { o.setRoot(nil); return o }
func (o TreeMap[K, V]) Clone() TreeMap[K, V] {
	clone := NewTreeMap[K, V]()
	o.ForEach(func(k K, v V) { clone.Set(k, v) })
	return clone
}

// ==============MapContainer接口=============
// 遍历和转换
func (o TreeMap[K, V]) ForEach(f func(k K, v V)) TreeMap[K, V] {
	entries := make([]*MapEntry[K, V], 0, o.Len())
	var _forEach func(node *avlNode[K, V])
	_forEach = func(node *avlNode[K, V]) {
		if node == nil {
			return
		}
		_forEach(node.left)
		entries = append(entries, newEntryFromNode(node))
		_forEach(node.right)
	}
	_forEach(o.root())
	for _, entry := range entries { // 支持边遍历边修改
		f(entry.K, entry.V)
	}
	return o
}
func (o TreeMap[K, V]) ToMap() map[K]V {
	m := make(map[K]V, o.Len())
	o.ForEach(func(k K, v V) { m[k] = v })
	return m
}

// 基本操作
func (o TreeMap[K, V]) Get(k K) V {
	node := o.getNode(k)
	if node != nil {
		return node.v
	}
	if o.factory != nil { // 使用工厂函数设置值
		v := o.factory()
		o.Set(k, v)
		return v
	}
	return *new(V)
}
func (o TreeMap[K, V]) Set(k K, v V) TreeMap[K, V] {
	var _set func(node *avlNode[K, V], k K, v V) *avlNode[K, V]
	_set = func(node *avlNode[K, V], k K, v V) *avlNode[K, V] {
		if node == nil {
			return newAVLNode(k, v).updateStatus(o.w)
		} else if k < node.k {
			node.left = _set(node.left, k, v)
		} else if k > node.k {
			node.right = _set(node.right, k, v)
		} else {
			node.v = v
		}
		return node.updateStatus(o.w).balance(o.w)
	}

	o.setRoot(_set(o.root(), k, v))
	return o
}
func (o TreeMap[K, V]) Del(k K) TreeMap[K, V] {
	var _del func(node *avlNode[K, V], k K) *avlNode[K, V]
	_del = func(node *avlNode[K, V], k K) *avlNode[K, V] {
		if node == nil {
			return nil
		}
		if k < node.k {
			node.left = _del(node.left, k)
		} else if k > node.k {
			node.right = _del(node.right, k)
		} else if node.right == nil { // 只有左子树或者右子树为零，才会真正发生删除
			return node.left
		} else if node.left == nil {
			return node.right
		} else {
			next := node.right
			for next.left != nil { // 找到后继节点
				next = next.left
			}
			node.k, node.v = next.k, next.v
			node.right = _del(node.right, next.k)
		}
		return node.updateStatus(o.w).balance(o.w)
	}

	o.setRoot(_del(o.root(), k))
	return o
}
func (o TreeMap[K, V]) Has(k K) bool { return o.getNode(k) != nil }

// 辅助函数，返回键对应的节点
func (o TreeMap[K, V]) getNode(k K) *avlNode[K, V] {
	node := o.root()
	for node != nil {
		if k < node.k {
			node = node.left
		} else if k > node.k {
			node = node.right
		} else {
			return node
		}
	}
	return nil
}

// 组合操作
func (o TreeMap[K, V]) GetOr(k K, defalutValue V) V {
	if node := o.getNode(k); node != nil {
		return node.v
	}
	return defalutValue
}
func (o TreeMap[K, V]) GetOrSet(k K, v V) V {
	if node := o.getNode(k); node != nil {
		return node.v
	}
	o.Set(k, v)
	return v
}
func (o TreeMap[K, V]) Extend(other TreeMap[K, V]) TreeMap[K, V] {
	other.ForEach(func(k K, v V) { o.Set(k, v) })
	return o
}
func (o TreeMap[K, V]) DelFunc(f func(K, V) bool) TreeMap[K, V] {
	return o.ForEach(func(k K, v V) {
		if f(k, v) {
			o.Del(k)
		}
	})
}
func (o TreeMap[K, V]) ReplaceFunc(f func(K, V) V) TreeMap[K, V] {
	return o.ForEach(func(k K, v V) { o.Set(k, f(k, v)) })
}

// 键值查询
func (o TreeMap[K, V]) Keys() []K {
	var keys []K
	o.ForEach(func(k K, v V) { keys = append(keys, k) })
	return keys
}
func (o TreeMap[K, V]) Values() []V {
	var values []V
	o.ForEach(func(k K, v V) { values = append(values, v) })
	return values
}

// ==============TreeMapContainer接口=============
// 二分查找键值对
func (o TreeMap[K, V]) First() *MapEntry[K, V] {
	if o.root() == nil {
		return nil
	}
	node := o.root()
	for node.left != nil {
		node = node.left
	}
	return newEntryFromNode(node)
}
func (o TreeMap[K, V]) Last() *MapEntry[K, V] {
	if o.root() == nil {
		return nil
	}
	node := o.root()
	for node.right != nil {
		node = node.right
	}
	return newEntryFromNode(node)
}
func (o TreeMap[K, V]) Lower(k K) *MapEntry[K, V] {
	var res *avlNode[K, V]
	var _lowerKey func(node *avlNode[K, V], k K)
	_lowerKey = func(node *avlNode[K, V], k K) {
		if node == nil {
			return
		}
		if node.k < k {
			res = node
			_lowerKey(node.right, k) // 寻找更大的键，满足<k
		} else {
			_lowerKey(node.left, k)
		}
	}

	_lowerKey(o.root(), k)
	return newEntryFromNode(res)
}
func (o TreeMap[K, V]) Higher(k K) *MapEntry[K, V] {
	var res *avlNode[K, V]
	var _higherKey func(node *avlNode[K, V], k K)
	_higherKey = func(node *avlNode[K, V], k K) {
		if node == nil {
			return
		}
		if node.k > k {
			res = node
			_higherKey(node.left, k) // 寻找更小的键，满足>k
		} else {
			_higherKey(node.right, k)
		}
	}

	_higherKey(o.root(), k)
	return newEntryFromNode(res)
}
func (o TreeMap[K, V]) Floor(k K) *MapEntry[K, V] {
	var res *avlNode[K, V]
	var _floorKey func(node *avlNode[K, V], k K)
	_floorKey = func(node *avlNode[K, V], k K) {
		if node == nil {
			return
		}
		if node.k <= k {
			res = node
			_floorKey(node.right, k) // 寻找更大的键，满足<=k
		} else {
			_floorKey(node.left, k)
		}
	}

	_floorKey(o.root(), k)
	return newEntryFromNode(res)
}
func (o TreeMap[K, V]) Ceiling(k K) *MapEntry[K, V] {
	var res *avlNode[K, V]
	var _ceilingKey func(node *avlNode[K, V], k K)
	_ceilingKey = func(node *avlNode[K, V], k K) {
		if node == nil {
			return
		}
		if node.k >= k {
			res = node
			_ceilingKey(node.left, k) // 寻找更小的键，满足>=k
		} else {
			_ceilingKey(node.right, k)
		}
	}

	_ceilingKey(o.root(), k)
	return newEntryFromNode(res)
}

// 排名相关操作
func (o TreeMap[K, V]) Select(i int) *MapEntry[K, V] {
	var _select func(node *avlNode[K, V], i int) *avlNode[K, V]
	_select = func(node *avlNode[K, V], i int) *avlNode[K, V] {
		if node == nil {
			return nil
		}
		lSize := o.size(node.left)
		if lSize == i-1 {
			return node
		} else if lSize >= i {
			return _select(node.left, i)
		} else {
			return _select(node.right, i-lSize-1)
		}
	}

	return newEntryFromNode(_select(o.root(), i))
}
func (o TreeMap[K, V]) Rank(k K) int { // 返回小于等于k的节点个数
	var _rank func(node *avlNode[K, V]) int
	_rank = func(node *avlNode[K, V]) int {
		if node == nil {
			return 0
		}
		if k == node.k {
			return o.size(node.left) + 1
		} else if k < node.k {
			return _rank(node.left) // 当小于所有节点，返回0
		} else {
			return o.size(node.left) + 1 + _rank(node.right) // 当大于所有节点，会返回n
		}
	}
	return _rank(o.root())
}

// 辅助函数，返回节点大小
func (o TreeMap[K, V]) size(node *avlNode[K, V]) int {
	if node == nil {
		return 0
	}
	return node.size
}

// 子映射截取
func (o TreeMap[K, V]) HeadMap(n int) TreeMap[K, V] {
	newMap := NewTreeMap[K, V]()
	var _headMap func(node *avlNode[K, V])
	_headMap = func(node *avlNode[K, V]) {
		if node == nil {
			return
		}
		_headMap(node.left)
		if n <= 0 {
			return
		}
		newMap.Set(node.k, node.v)
		n--
		_headMap(node.right)
	}

	_headMap(o.root())
	return newMap
}
func (o TreeMap[K, V]) TailMap(n int) TreeMap[K, V] {
	newMap := NewTreeMap[K, V]()
	var _tailMap func(node *avlNode[K, V])
	_tailMap = func(node *avlNode[K, V]) {
		if node == nil {
			return
		}
		_tailMap(node.right)
		if n <= 0 {
			return
		}
		newMap.Set(node.k, node.v)
		n--
		_tailMap(node.left)
	}

	_tailMap(o.root())
	return newMap
}

// ==============其他============= */
func (o TreeMap[K, V]) PrintTree() {
	var _print func(node *avlNode[K, V], depth int)
	_print = func(node *avlNode[K, V], depth int) {
		if node == nil {
			return
		}
		_print(node.right, depth+1)
		fmt.Printf("%v[%v,%v]\n", strings.Repeat("  ", depth*2), node.k, node.v)
		_print(node.left, depth+1)
	}
	_print(o.root(), 0)
	fmt.Println()
}

// ===================================计数Map===================================
// TOOD
// type TreeMapCounter struct {
// 	*TreeMap[K, int]
// }
//
// func NewTreeMapCounter[K cmp.Ordered]() *TreeMapCounter {
// 	return &TreeMapCounter{NewDefaultTreeMap[K, int](func() int { return 0 })}
// }

// ===================================有序集===================================

type TreeSet[T cmp.Ordered] struct {
	m TreeMap[T, struct{}]
}

func NewTreeSet[T cmp.Ordered]() TreeSet[T] {
	return TreeSet[T]{m: NewTreeMap[T, struct{}]()}
}

func NewTreeSetFromSlice[T cmp.Ordered](slice []T) TreeSet[T] {
	o := NewTreeSet[T]()
	for _, v := range slice {
		o.Add(v)
	}
	return o
}

// ==============ValueContainer接口=============
func (o TreeSet[T]) Len() int          { return o.m.Len() }
func (o TreeSet[T]) String() string    { return "TreeSet" + fmt.Sprint(o.ToSlice()) }
func (o TreeSet[T]) Clear() TreeSet[T] { o.m.Clear(); return o }
func (o TreeSet[T]) Clone() TreeSet[T] { o.m = o.m.Clone(); return o }
func (o TreeSet[T]) ForEach(f func(x T)) TreeSet[T] {
	o.m.ForEach(func(k T, v struct{}) { f(k) })
	return o
}
func (o TreeSet[T]) ToSlice() []T { return o.m.Keys() }

// ==============SetContainer接口=============
// 基本操作
func (o TreeSet[T]) Add(x T) TreeSet[T] { o.m.Set(x, struct{}{}); return o }
func (o TreeSet[T]) Del(x T) TreeSet[T] { o.m.Del(x); return o }
func (o TreeSet[T]) Has(x T) bool       { return o.m.Has(x) }

// 组合操作
func (o TreeSet[T]) DelFunc(f func(x T) bool) TreeSet[T] {
	return o.ForEach(func(x T) {
		if f(x) {
			o.Del(x)
		}
	})
}
func (o TreeSet[T]) HasSubset(other TreeSet[T]) bool {
	for _, x := range other.ToSlice() {
		if !o.Has(x) {
			return false
		}
	}
	return true
}
func (o TreeSet[T]) Union(other TreeSet[T]) TreeSet[T] {
	res := NewTreeSet[T]()
	o.ForEach(func(k T) { res.Add(k) })
	other.ForEach(func(k T) { res.Add(k) })
	return res
}
func (o TreeSet[T]) Intersect(other TreeSet[T]) TreeSet[T] {
	res := NewTreeSet[T]()
	o.ForEach(func(k T) {
		if other.Has(k) {
			res.Add(k)
		}
	})
	return res
}
func (o TreeSet[T]) Difference(other TreeSet[T]) TreeSet[T] {
	res := NewTreeSet[T]()
	o.ForEach(func(k T) {
		if !other.Has(k) {
			res.Add(k)
		}
	})
	return res
}

// ==============TreeSetContainer接口=============
// 子集合截取
func (o TreeSet[T]) HeadSet(x T) TreeSet[T] {
	return NewTreeSetFromSlice(o.m.HeadMap(o.m.Rank(x)).Keys())
}
func (o TreeSet[T]) TailSet(x T) TreeSet[T] {
	return NewTreeSetFromSlice(o.m.TailMap(o.m.Rank(x)).Keys())
}

// 排名相关的查找
func (o TreeSet[T]) Rank(x T) int { return o.m.Rank(x) }
func (o TreeSet[T]) Kth(k int) T  { return o.m.Select(k).K }

// 二分查找值，可能会panic
func (o TreeSet[T]) First() T      { return o.m.First().K }
func (o TreeSet[T]) Last() T       { return o.m.Last().K }
func (o TreeSet[T]) Lower(x T) T   { return o.m.Lower(x).K }
func (o TreeSet[T]) Higher(x T) T  { return o.m.Higher(x).K }
func (o TreeSet[T]) Floor(x T) T   { return o.m.Floor(x).K }
func (o TreeSet[T]) Ceiling(x T) T { return o.m.Ceiling(x).K }

// ===================================多重有序集===================================

type MultiTreeSet[T cmp.Ordered] struct {
	m      TreeMap[T, int]
	sumMap map[T]int
}

func NewMultiTreeSet[T cmp.Ordered]() MultiTreeSet[T] {
	sumMap := make(map[T]int)
	return MultiTreeSet[T]{
		sumMap: sumMap,
		m: NewTreeMap[T, int]().WithFactory(func() int { return 0 }).WithWatcher(func(cur, left, right *MapEntry[T, int]) {
			lCnt, rCnt := 0, 0
			if left != nil {
				lCnt = sumMap[left.K]
			}
			if right != nil {
				rCnt = sumMap[right.K]
			}
			sumMap[cur.K] = cur.V + lCnt + rCnt
		})}
}

func NewMultiTreeSetFromSlice[T cmp.Ordered](slice []T) MultiTreeSet[T] {
	o := NewMultiTreeSet[T]()
	for _, v := range slice {
		o.Add(v)
	}
	return o
}

func NewMultiTreeSetFromMap[T cmp.Ordered](m map[T]int) MultiTreeSet[T] {
	o := NewMultiTreeSet[T]()
	for k, v := range m {
		o.AddN(k, v)
	}
	return o
}

// --------------------其他转换--------------------
func (o MultiTreeSet[T]) ToTreeSet() TreeSet[T] { return NewTreeSetFromSlice(o.ToSlice()) } // 去重

// --------------------ValueContainer接口--------------------
func (o MultiTreeSet[T]) Len() int { return o.m.Len() }
func (o MultiTreeSet[T]) Clear() MultiTreeSet[T] {
	o.m.Clear()
	clear(o.sumMap)
	return o
}
func (o MultiTreeSet[T]) Clone() MultiTreeSet[T] {
	o.m = o.m.Clone()
	o.sumMap = maps.Clone(o.sumMap)
	return o
}
func (o MultiTreeSet[T]) ForEach(f func(x T)) { o.ForEachCnt(func(k T, v int) { f(k) }) }
func (o MultiTreeSet[T]) ToSlice() []T {
	res := make([]T, 0, o.Total())
	o.ForEachCnt(func(k T, cnt int) {
		for range cnt {
			res = append(res, k)
		}
	})
	return res
}
func (o MultiTreeSet[T]) String() string {
	return strings.ReplaceAll(o.m.String(), "TreeMap", "MultiTreeSet")
}

// --------------------Set接口--------------------

// 基本操作
func (o MultiTreeSet[T]) Has(x T) bool            { return o.m.Has(x) }
func (o MultiTreeSet[T]) Del(x T) MultiTreeSet[T] { return o.DelN(x, 1) }
func (o MultiTreeSet[T]) Add(x T) MultiTreeSet[T] { return o.AddN(x, 1) }

// 组合操作
func (o MultiTreeSet[T]) DelFunc(f func(x T) bool) MultiTreeSet[T] {
	return o.ForEachCnt(func(x T, cnt int) {
		if f(x) {
			o.Del(x)
		}
	})
}
func (o MultiTreeSet[T]) HasSubset(other MultiTreeSet[T]) bool {
	for v, cnt := range other.ToMap() {
		if cnt > o.Count(v) {
			return false
		}
	}
	return true
}

func (o MultiTreeSet[T]) Intersect(other MultiTreeSet[T]) MultiTreeSet[T] {
	res := NewMultiTreeSet[T]()
	other.ForEachCnt(func(v T, cnt int) { res.AddN(v, min(cnt, o.Count(v))) })
	return res
}
func (o MultiTreeSet[T]) Union(other MultiTreeSet[T]) MultiTreeSet[T] {
	res := o.Clone()
	other.ForEachCnt(func(v T, cnt int) { res.AddN(v, cnt-o.Count(v)) })
	return res
}
func (o MultiTreeSet[T]) Difference(other MultiTreeSet[T]) MultiTreeSet[T] {
	res := o.Clone()
	other.ForEachCnt(func(v T, cnt int) { res.DelN(v, o.Count(v)-cnt) })
	return res
}

// --------------------MultiSet接口--------------------

// 遍历和转换
func (o MultiTreeSet[T]) ForEachCnt(f func(x T, cnt int)) MultiTreeSet[T] { o.m.ForEach(f); return o }
func (o MultiTreeSet[T]) ToMap() map[T]int {
	res := make(map[T]int, o.Len())
	o.ForEachCnt(func(k T, cnt int) { res[k] = cnt })
	return res
}

// 计数操作
func (o MultiTreeSet[T]) AddN(x T, n int) MultiTreeSet[T] {
	if n <= 0 {
		return o
	}
	o.m.Set(x, o.m.Get(x)+n)
	return o
}
func (o MultiTreeSet[T]) DelN(x T, n int) MultiTreeSet[T] {
	if n <= 0 || !o.m.Has(x) {
		return o
	}
	if n >= o.m.Get(x) {
		o.m.Del(x)
	} else {
		o.m.Set(x, o.m.Get(x)-n)
	}
	return o
}
func (o MultiTreeSet[T]) DelAll(x T) MultiTreeSet[T] { return o.DelN(x, o.Count(x)) }
func (o MultiTreeSet[T]) Count(x T) int {
	if !o.m.Has(x) {
		return 0
	}
	return o.m.Get(x)
}
func (o MultiTreeSet[T]) Total() int {
	if o.Len() == 0 {
		return 0
	}
	return o.sumMap[o.m.root().k]
}
func (o MultiTreeSet[T]) DelAllFunc(f func(T, int) bool) MultiTreeSet[T] {
	return o.ForEachCnt(func(x T, cnt int) {
		if f(x, cnt) {
			o.DelAll(x)
		}
	})
}
func (o MultiTreeSet[T]) ReplaceFunc(f func(T, int) int) MultiTreeSet[T] {
	return o.ForEachCnt(func(x T, cnt int) {
		newCnt := f(x, cnt)
		if newCnt > cnt {
			o.AddN(x, newCnt-cnt)
		} else if newCnt < cnt {
			o.DelN(x, cnt-newCnt)
		}
	})
}

// --------------------TreeSet接口--------------------

// 二分查找值，可能会panic
func (o MultiTreeSet[T]) First() T      { return o.m.First().K }
func (o MultiTreeSet[T]) Last() T       { return o.m.Last().K }
func (o MultiTreeSet[T]) Lower(x T) T   { return o.m.Lower(x).K }
func (o MultiTreeSet[T]) Higher(x T) T  { return o.m.Higher(x).K }
func (o MultiTreeSet[T]) Floor(x T) T   { return o.m.Floor(x).K }
func (o MultiTreeSet[T]) Ceiling(x T) T { return o.m.Ceiling(x).K }

// 截取子集合（不考虑重数）
func (o MultiTreeSet[T]) HeadSet(k int) MultiTreeSet[T] {
	return NewMultiTreeSetFromMap(o.m.HeadMap(k).ToMap())
}
func (o MultiTreeSet[T]) TailSet(k int) MultiTreeSet[T] {
	return NewMultiTreeSetFromMap(o.m.TailMap(k).ToMap())
}

// 排名相关的查找（考虑重数）
func (o MultiTreeSet[T]) Rank(x T) int { // 取最小排名
	var rank func(node *MapEntry[T, int]) int
	rank = func(node *MapEntry[T, int]) int {
		if node == nil {
			return 0
		}
		if x == node.K {
			return o.sum(node.Left()) + 1
		} else if x < node.K {
			return rank(node.Left())
		} else {
			return o.sum(node.Left()) + o.Count(node.K) + rank(node.Right())
		}
	}
	return rank(o.m.GetRoot())
}
func (o MultiTreeSet[T]) Select(i int) T { // 如果第i个元素不存在会panic
	var _select func(node *MapEntry[T, int], i int) *MapEntry[T, int]
	_select = func(node *MapEntry[T, int], i int) *MapEntry[T, int] {
		if node == nil {
			return nil
		}
		lCnt := o.sum(node.Left())
		curCnt := o.Count(node.K)
		if lCnt+1 <= i && i <= lCnt+curCnt {
			return node
		} else if i <= lCnt {
			return _select(node.Left(), i)
		} else {
			return _select(node.Right(), i-lCnt-curCnt)
		}
	}
	return _select(o.m.GetRoot(), i).K
}

// 辅助函数，计算节点对应子树的计数
func (o MultiTreeSet[T]) sum(node *MapEntry[T, int]) int {
	if node == nil {
		return 0
	}
	return o.sumMap[node.K]
}
