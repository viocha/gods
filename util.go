package gods

import (
	"fmt"
	"reflect"
	"strings"
)

// ===========================存储可能不存在的值的结构体============================

type Optional[T any] struct{ v *T }

// 创建一个空的Optional对象
func NewOptional[T any]() Optional[T] {
	return Optional[T]{}
}

// 设置指定的值，返回本身
func (o Optional[T]) WithValue(v T) Optional[T] { o.v = &v; return o }

// 获取值，如果不存在会panic
func (o Optional[T]) Get() T {
	if !o.Exists() {
		panic("Optional不含有值")
	}
	return *o.v
}

func (o Optional[T]) GetOr(defalutValue T) T {
	if !o.Exists() {
		return defalutValue
	}
	return *o.v
}

func (o Optional[T]) Exists() bool {
	return o.v != nil
}

// ===================================树状打印工具===================================

// 模拟二叉树节点的接口，左右孩子要求返回指针
type TreeNodeLike interface {
	GetLeft() TreeNodeLike
	GetRight() TreeNodeLike
	NodeString() string
}

// 计算树的宽度和高度
func getWidAndHeight(root TreeNodeLike) (w, h int) {
	if IsNil(root) {
		return 0, 0
	}
	nodeWidth := len(root.NodeString())
	lw, lh := getWidAndHeight(root.GetLeft())
	rw, rh := getWidAndHeight(root.GetRight())

	if lw == 0 && rw == 0 {
		return nodeWidth, 1
	}
	return lw + rw + nodeWidth, max(lh, rh) + 2 // 包括引导线
}

// --------------------打印树状结构的画板--------------------

type DrawBoard struct {
	oi, oj int
	board  [][]rune
}

// 根据树的宽度和高度，创建一个画板
func createDrawBoard(root TreeNodeLike) *DrawBoard {
	w, h := getWidAndHeight(root)
	board := make([][]rune, h)
	for i := range board {
		board[i] = make([]rune, w)
		for j := range board[i] {
			board[i][j] = ' '
		}
	}
	return &DrawBoard{board: board}
}

// 设置画板的原点
func (d *DrawBoard) setOrigin(i, j int) *DrawBoard {
	d.oi, d.oj = i, j
	return d
}

// 在相对原点的指定位置绘制字符串，限制在j方向上，i方向等于0
func (d *DrawBoard) drawY(j int, s string) *DrawBoard {
	for k, b := range []rune(s) {
		d.board[d.oi+0][d.oj+j+k] = b
	}
	return d
}

// 将画板的图像转换成字符串
func (d *DrawBoard) String() string {
	sb := &strings.Builder{}
	for _, line := range d.board {
		sb.WriteString(string(line))
		sb.WriteString("\n")
	}
	return sb.String()
}

// --------------------递归绘制树形结构--------------------

func IsNil(v any) bool { return v == nil || reflect.ValueOf(v).IsNil() }

// 绘制一个子树，返回子树的宽度，以及子树根节点的左侧偏移
func drawRoot(d *DrawBoard, oi, oj int, root TreeNodeLike) (width, leftWidth int) {
	if IsNil(root) {
		return
	}
	nodeStr := root.NodeString()
	nodeWidth := len([]rune(nodeStr))
	lw, llw := drawRoot(d, oi+2, oj, root.GetLeft())
	rw, rlw := drawRoot(d, oi+2, oj+lw+nodeWidth, root.GetRight())

	d.setOrigin(oi, oj)
	d.drawY(lw, nodeStr)
	d.setOrigin(oi+1, oj)
	drawGuide(d, nodeWidth, lw, llw, rw, rlw)

	width = lw + nodeWidth + rw
	leftWidth = lw + nodeWidth/2
	return width, leftWidth
}

// 绘制引导线： ┏━┻━┓  ━┛ ┗━
func drawGuide(d *DrawBoard, nodeWidth int, lw int, llw int, rw int, rlw int) {
	drawLine := func(l, r int) {
		for j := l; j <= r; j++ {
			d.drawY(j, "━")
		}
	}
	if lw == 0 && rw != 0 { // 无左子树
		m := nodeWidth / 2
		r := nodeWidth + rlw
		d.drawY(m, "┗")
		drawLine(m+1, r-1)
		d.drawY(r, "┓")
	} else if lw != 0 && rw == 0 { // 无右子树
		l := llw
		m := lw + nodeWidth/2
		d.drawY(l, "┏")
		drawLine(l+1, m-1)
		d.drawY(m, "┛")
	} else if lw != 0 && rw != 0 { // 有左右子树
		l := llw
		m := lw + nodeWidth/2
		r := lw + nodeWidth + rlw
		d.drawY(l, "┏")
		drawLine(l+1, m-1)
		d.drawY(m, "┻")
		drawLine(m+1, r-1)
		d.drawY(r, "┓")
	}
}

// 打印树状结构
func PrintTree(root TreeNodeLike) {
	d := createDrawBoard(root)
	drawRoot(d, 0, 0, root)
	fmt.Println(d.String())
}
