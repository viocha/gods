package gods

import (
	"fmt"
	"reflect"
	"testing"
)

// ======================测试树节点============================
type TreeNode struct {
	v    int
	l, r *TreeNode
}

func (o TreeNode) GetLeft() TreeNodeLike { return o.l }

func (o TreeNode) GetRight() TreeNodeLike { return o.r }

func (o TreeNode) NodeString() string { return fmt.Sprint(o.v) }

func Test打印树结构(t *testing.T) {
	root := &TreeNode{v: 1}
	root.l = &TreeNode{v: 2}
	root.r = &TreeNode{v: 3}
	root.l.l = &TreeNode{v: 400}
	root.l.r = &TreeNode{v: 5}
	root.r.l = &TreeNode{v: 88}
	root.r.r = &TreeNode{v: 10000}

	PrintTree(root)
}

func TestNil(t *testing.T) {
	var nilPtr *int
	var anyVal any = nilPtr
	fmt.Println(anyVal, anyVal == nil) // <nil> false，动态类型是*int，动态值是nil，接口本身不为nil

	fmt.Println(reflect.ValueOf(anyVal).IsNil()) // true，判断动态值是否是nil
	//fmt.Println(reflect.ValueOf(nil).IsNil()) // 会报错，无法判断对nil调用IsNil()方法
}
