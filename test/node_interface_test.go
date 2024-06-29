package test

import (
	"fmt"
	"testing"
)

type BinaryTreeNode interface {
	GetValue() string
	GetLeft() BinaryTreeNode
	GetRight() BinaryTreeNode
}

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

func (n *Node) GetValue() string {
	return fmt.Sprint(n.Value)
}

func (n *Node) GetLeft() BinaryTreeNode {
	return n.Left
}

func (n *Node) GetRight() BinaryTreeNode {
	return n.Right
}

func TestNode(t *testing.T) {
	var root BinaryTreeNode = &Node{
		Value: 1,
		Left: &Node{
			Value: 2,
			Left:  &Node{Value: 4},
			Right: &Node{Value: 5},
		},
		Right: &Node{
			Value: 3,
			Left:  &Node{Value: 6},
			Right: &Node{Value: 7},
		},
	}

	// 使用接口方法
	fmt.Println(root.GetValue())                     // 1
	fmt.Println(root.GetLeft().GetValue())           // 2
	fmt.Println(root.GetRight().GetValue())          // 3
	fmt.Println(root.GetLeft().GetLeft().GetValue()) // 4

	// 验证接口方法的递归性
	leftChild := root.GetLeft()
	fmt.Println(leftChild.GetRight().GetValue()) // 5
}
