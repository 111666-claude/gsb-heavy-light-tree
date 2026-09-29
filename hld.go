// Package heavylight 用树链剖分支持路径求和与子树求和。
// 位置按「重儿子优先」的迭代 DFS 序分配：每条链连续，
// 每棵子树是区间 [pos, pos+size)。全程迭代，深链不爆栈。
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
	head   []int
	size   []int
	pos    []int
	prefix []int64
}

// New 建树。edges 里每条是 (parent, child)。
func New(count int, values []int64, edges [][2]int) (*Tree, error) {
	if count <= 0 || len(values) != count {
		return nil, ErrInput
	}
	children := make([][]int, count)
	parent := make([]int, count)
	for index := range parent {
		parent[index] = -1
	}
	for _, edge := range edges {
		from, to := edge[0], edge[1]
		if from < 0 || from >= count || to <= 0 || to >= count {
			return nil, ErrInput
		}
		if parent[to] != -1 {
			return nil, ErrInput
		}
		parent[to] = from
		children[from] = append(children[from], to)
	}

	// 迭代 DFS：求深度，并确认所有节点都从 0 可达（无环、无森林）。
	depth := make([]int, count)
	order := make([]int, 0, count)
	stack := []int{0}
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, node)
		for _, child := range children[node] {
			depth[child] = depth[node] + 1
			stack = append(stack, child)
		}
	}
	if len(order) != count {
		return nil, ErrInput
	}

	// 逆 DFS 序求子树 size 与重儿子（size 最大的孩子）。
	size := make([]int, count)
	heavy := make([]int, count)
	for index := range heavy {
		heavy[index] = -1
	}
	for index := len(order) - 1; index >= 0; index-- {
		node := order[index]
		size[node] = 1
		best, bestSize := -1, 0
		for _, child := range children[node] {
			size[node] += size[child]
			if size[child] > bestSize {
				best, bestSize = child, size[child]
			}
		}
		heavy[node] = best
	}

	// 重儿子优先的迭代位置分配：沿重链向下走，轻儿子压栈稍后处理。
	// 链内位置连续；子树所有后代先于栈中更早的节点被处理，故子树连续。
	pos := make([]int, count)
	head := make([]int, count)
	type chainRoot struct{ node, top int }
	todo := []chainRoot{{0, 0}}
	next := 0
	for len(todo) > 0 {
		root := todo[len(todo)-1]
		todo = todo[:len(todo)-1]
		for node := root.node; node != -1; node = heavy[node] {
			head[node] = root.top
			pos[node] = next
			next++
			for _, child := range children[node] {
				if child != heavy[node] {
					todo = append(todo, chainRoot{child, child})
				}
			}
		}
	}

	// 按位置重排点权并做前缀和。
	prefix := make([]int64, count+1)
	for node := 0; node < count; node++ {
		prefix[pos[node]+1] = values[node]
	}
	for index := 0; index < count; index++ {
		prefix[index+1] += prefix[index]
	}

	return &Tree{
		count:  count,
		value:  append([]int64(nil), values...),
		parent: parent,
		depth:  depth,
		heavy:  heavy,
		head:   head,
		size:   size,
		pos:    pos,
		prefix: prefix,
	}, nil
}

// rangeSum 返回位置区间 [left, right) 的权值和。
func (t *Tree) rangeSum(left, right int) int64 {
	return t.prefix[right] - t.prefix[left]
}

// SubtreeSum 返回以 node 为根的子树点权和，即区间 [pos, pos+size)。
func (t *Tree) SubtreeSum(node int) (int64, error) {
	if node < 0 || node >= t.count {
		return 0, ErrInput
	}
	return t.rangeSum(t.pos[node], t.pos[node]+t.size[node]), nil
}

// PathSum 返回 first 到 second 简单路径上的点权和。
// 沿两条链向上跳，链内用区间和，最后补上 LCA 所在那一段。
func (t *Tree) PathSum(first, second int) (int64, error) {
	if first < 0 || first >= t.count || second < 0 || second >= t.count {
		return 0, ErrInput
	}
	total := int64(0)
	for t.head[first] != t.head[second] {
		if t.depth[t.head[first]] < t.depth[t.head[second]] {
			first, second = second, first
		}
		top := t.head[first]
		total += t.rangeSum(t.pos[top], t.pos[first]+1)
		first = t.parent[top]
	}
	if t.depth[first] > t.depth[second] {
		first, second = second, first
	}
	total += t.rangeSum(t.pos[first], t.pos[second]+1)
	return total, nil
}
