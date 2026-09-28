// Package heavylight 用树链剖分支持路径求和与子树求和。
// 位置按「重儿子优先」的 DFS 序分配：每条链连续，每棵子树是区间 [pos, pos+size)。
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
	tree := &Tree{
		count:  count,
		value:  append([]int64(nil), values...),
		parent: make([]int, count),
		depth:  make([]int, count),
		heavy:  make([]int, count),
		head:   make([]int, count),
		size:   make([]int, count),
		pos:    make([]int, count),
		prefix: make([]int64, count+1),
	}
	for index := range tree.parent {
		tree.parent[index] = -1
		tree.heavy[index] = -1
	}
	for _, edge := range edges {
		parent, child := edge[0], edge[1]
		if parent < 0 || parent >= count || child <= 0 || child >= count {
			return nil, ErrInput
		}
		if tree.parent[child] != -1 {
			return nil, ErrInput
		}
		children[parent] = append(children[parent], child)
		tree.parent[child] = parent
	}

	// 迭代 DFS：记录访问顺序与深度（深链不能递归）。
	visit := make([]int, 0, count)
	stack := make([]int, 0, count)
	stack = append(stack, 0)
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		visit = append(visit, node)
		for _, child := range children[node] {
			tree.depth[child] = tree.depth[node] + 1
			stack = append(stack, child)
		}
	}
	if len(visit) != count {
		return nil, ErrInput
	}

	// 逆序求子树 size，并选 size 最大的孩子作重儿子。
	for index := len(visit) - 1; index >= 0; index-- {
		node := visit[index]
		tree.size[node] = 1
		best, bestSize := -1, 0
		for _, child := range children[node] {
			tree.size[node] += tree.size[child]
			if tree.size[child] > bestSize {
				best, bestSize = child, tree.size[child]
			}
		}
		tree.heavy[node] = best
	}

	// 重儿子优先的迭代 DFS 位置分配：沿重链连续编号，轻儿子入栈稍后处理。
	// 栈是 LIFO，处理某条链期间入栈的轻儿子（即链上节点的其余后代）
	// 会先于更早入栈的节点出栈，因此每棵子树恰好占连续区间。
	type frame struct{ node, head int }
	frames := make([]frame, 0, count)
	frames = append(frames, frame{0, 0})
	next := 0
	for len(frames) > 0 {
		top := frames[len(frames)-1]
		frames = frames[:len(frames)-1]
		for node := top.node; node != -1; node = tree.heavy[node] {
			tree.head[node] = top.head
			tree.pos[node] = next
			next++
			for _, child := range children[node] {
				if child != tree.heavy[node] {
					frames = append(frames, frame{child, child})
				}
			}
		}
	}

	byPos := make([]int64, count)
	for node := 0; node < count; node++ {
		byPos[tree.pos[node]] = tree.value[node]
	}
	for index := 0; index < count; index++ {
		tree.prefix[index+1] = tree.prefix[index] + byPos[index]
	}
	return tree, nil
}

// SubtreeSum 返回以 node 为根的子树点权和，即区间 [pos, pos+size) 的和。
func (t *Tree) SubtreeSum(node int) (int64, error) {
	if node < 0 || node >= t.count {
		return 0, ErrInput
	}
	left := t.pos[node]
	right := left + t.size[node]
	return t.prefix[right] - t.prefix[left], nil
}

// PathSum 返回 first 到 second 简单路径上的点权和。
// 沿两条链向上跳，链内用前缀和求区间和，最后补上 LCA 所在链的一段。
func (t *Tree) PathSum(first, second int) (int64, error) {
	if first < 0 || second < 0 || first >= t.count || second >= t.count {
		return 0, ErrInput
	}
	total := int64(0)
	for t.head[first] != t.head[second] {
		if t.depth[t.head[first]] < t.depth[t.head[second]] {
			first, second = second, first
		}
		top := t.head[first]
		total += t.prefix[t.pos[first]+1] - t.prefix[t.pos[top]]
		first = t.parent[top]
	}
	if t.depth[first] > t.depth[second] {
		first, second = second, first
	}
	total += t.prefix[t.pos[second]+1] - t.prefix[t.pos[first]]
	return total, nil
}
