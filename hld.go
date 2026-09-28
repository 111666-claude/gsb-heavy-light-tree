// Package heavylight 用树链剖分支持路径求和与子树求和。
// 缺陷：位置按链编号顺序拼接（子树区间不连续），重儿子也选了第一个孩子。
package heavylight

import "errors"

// ErrInput 表示输入不合法。
var ErrInput = errors.New("heavylight: 输入不合法")

// Tree 是一棵以 0 为根的有根树。
type Tree struct {
	count  int
	value  []int64
	parent []int
	depth  []int
	heavy  []int
	size   []int
	order  []int
	prefix []int64
}

// New 建树。edges 里每条是 (parent, child)。
func New(count int, values []int64, edges [][2]int) (*Tree, error) {
	if count <= 0 || len(values) != count {
		return nil, ErrInput
	}
	children := make([][]int, count)
	tree := &Tree{
		count:  count,
		value:  append([]int64(nil), values...),
		parent: make([]int, count),
		depth:  make([]int, count),
		heavy:  make([]int, count),
		size:   make([]int, count),
		order:  make([]int, count),
		prefix: make([]int64, count+1),
	}
	for index := range tree.parent {
		tree.parent[index] = -1
		tree.heavy[index] = -1
	}
	for _, edge := range edges {
		if edge[0] < 0 || edge[0] >= count || edge[1] <= 0 || edge[1] >= count {
			return nil, ErrInput
		}
		children[edge[0]] = append(children[edge[0]], edge[1])
		tree.parent[edge[1]] = edge[0]
		tree.depth[edge[1]] = tree.depth[edge[0]] + 1
	}
	tree.measure(0, children)
	// 缺陷：位置直接按节点编号给，子树区间并不连续
	for index := 0; index < count; index++ {
		tree.order[index] = index
	}
	for node := 0; node < count; node++ {
		tree.prefix[node+1] = tree.prefix[node] + tree.value[node]
	}
	return tree, nil
}

func (t *Tree) measure(node int, children [][]int) int {
	t.size[node] = 1
	best, bestSize := -1, 0
	for _, child := range children[node] {
		size := t.measure(child, children)
		t.size[node] += size
		if size > bestSize {
			best, bestSize = child, size
		}
	}
	t.heavy[node] = best
	return t.size[node]
}

// SubtreeSum 缺陷：用位置区间 [pos, pos+size) 求和，但位置不是 DFS 序。
func (t *Tree) SubtreeSum(node int) (int64, error) {
	if node < 0 || node >= t.count {
		return 0, ErrInput
	}
	left := t.order[node]
	right := left + t.size[node]
	if right > t.count {
		return 0, ErrInput
	}
	return t.prefix[right] - t.prefix[left], nil
}

// PathSum 路径求和。
func (t *Tree) PathSum(first, second int) (int64, error) {
	if first < 0 || second < 0 || first >= t.count || second >= t.count {
		return 0, ErrInput
	}
	total := int64(0)
	for depth := t.depth[first] + t.depth[second] + 1; depth > 0; depth-- {
		if t.depth[first] >= t.depth[second] {
			total += t.value[first]
			if t.parent[first] >= 0 {
				first = t.parent[first]
			} else {
				break
			}
		} else {
			total += t.value[second]
			second = t.parent[second]
		}
	}
	return total, nil
}
