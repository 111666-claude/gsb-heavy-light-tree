package heavylight

import (
	"math/rand"
	"testing"
)

func chain(count int) [][2]int {
	edges := [][2]int{}
	for index := 1; index < count; index++ {
		edges = append(edges, [2]int{index - 1, index})
	}
	return edges
}

func TestBadInput(t *testing.T) {
	if _, err := New(0, nil, nil); err != ErrInput {
		t.Fatalf("空树应该报错：%v", err)
	}
	if _, err := New(2, []int64{1}, nil); err != ErrInput {
		t.Fatalf("值个数不对应该报错：%v", err)
	}
}

func TestSingleNode(t *testing.T) {
	tree, err := New(1, []int64{7}, nil)
	if err != nil {
		t.Fatalf("New 出错：%v", err)
	}
	if sum, _ := tree.SubtreeSum(0); sum != 7 {
		t.Fatalf("单点子树和应该是 7：%d", sum)
	}
}

func TestPathOnChain(t *testing.T) {
	tree, err := New(4, []int64{1, 2, 3, 4}, chain(4))
	if err != nil {
		t.Fatalf("New 出错：%v", err)
	}
	if sum, _ := tree.PathSum(0, 3); sum != 10 {
		t.Fatalf("0..3 路径和应该是 10：%d", sum)
	}
}

func TestSubtreeOfRoot(t *testing.T) {
	tree, err := New(4, []int64{1, 2, 3, 4}, chain(4))
	if err != nil {
		t.Fatalf("New 出错：%v", err)
	}
	if sum, _ := tree.SubtreeSum(0); sum != 10 {
		t.Fatalf("根子树和应该是 10：%d", sum)
	}
}

func TestSize(t *testing.T) {
	tree, _ := New(3, []int64{1, 1, 1}, [][2]int{{0, 1}, {1, 2}})
	if tree.count != 3 {
		t.Fatalf("节点数不对：%d", tree.count)
	}
}

func brutePathSum(values []int64, parent []int, depth []int, first, second int) int64 {
	total := int64(0)
	for first != second {
		if depth[first] < depth[second] {
			first, second = second, first
		}
		total += values[first]
		first = parent[first]
	}
	return total + values[first]
}

func bruteSubtreeSum(values []int64, children [][]int, root int) int64 {
	total := int64(0)
	stack := []int{root}
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		total += values[node]
		stack = append(stack, children[node]...)
	}
	return total
}

func buildBrute(count int, edges [][2]int) ([]int, []int, [][]int) {
	parent := make([]int, count)
	depth := make([]int, count)
	children := make([][]int, count)
	parent[0] = -1
	for _, edge := range edges {
		parent[edge[1]] = edge[0]
		children[edge[0]] = append(children[edge[0]], edge[1])
	}
	order := []int{0}
	for head := 0; head < len(order); head++ {
		node := order[head]
		for _, child := range children[node] {
			depth[child] = depth[node] + 1
			order = append(order, child)
		}
	}
	return parent, depth, children
}

func TestHeavyChoice(t *testing.T) {
	// 0 的孩子里 2 的子树更大，重儿子必须是 2 而不是先出现的 1。
	edges := [][2]int{{0, 1}, {0, 2}, {2, 3}, {2, 4}}
	tree, err := New(5, []int64{1, 1, 1, 1, 1}, edges)
	if err != nil {
		t.Fatalf("New 出错：%v", err)
	}
	if tree.heavy[0] != 2 {
		t.Fatalf("0 的重儿子应该是 2：%d", tree.heavy[0])
	}
	if tree.heavy[2] != 3 && tree.heavy[2] != 4 {
		t.Fatalf("2 的重儿子应该是 3 或 4：%d", tree.heavy[2])
	}
	if tree.heavy[1] != -1 || tree.heavy[3] != -1 {
		t.Fatalf("叶子不该有重儿子")
	}
	// 链归属：0、2、3 在同一条链，1 自成一链。
	if tree.head[0] != 0 || tree.head[2] != 0 || tree.head[3] != 0 {
		t.Fatalf("重链链头应该都是 0：%v", tree.head)
	}
	if tree.head[1] != 1 || tree.head[4] != 4 {
		t.Fatalf("轻儿子应该自成链头：%v", tree.head)
	}
}

func TestSubtreeInterval(t *testing.T) {
	// 题目用例：subtree(1) = {1,3} = 6+8 = 14。
	tree, err := New(4, []int64{5, 6, 7, 8}, [][2]int{{0, 1}, {0, 2}, {1, 3}})
	if err != nil {
		t.Fatalf("New 出错：%v", err)
	}
	if sum, _ := tree.SubtreeSum(1); sum != 14 {
		t.Fatalf("subtree(1) 应该是 14：%d", sum)
	}
	// 每个点的子树位置集合恰好是区间 [pos, pos+size)。
	count := 4
	_, _, children := buildBrute(count, [][2]int{{0, 1}, {0, 2}, {1, 3}})
	for node := 0; node < count; node++ {
		inside := make([]bool, count)
		stack := []int{node}
		for len(stack) > 0 {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			pos := tree.pos[cur]
			if pos < tree.pos[node] || pos >= tree.pos[node]+tree.size[node] {
				t.Fatalf("节点 %d 的位置 %d 不在 %d 的子树区间内", cur, pos, node)
			}
			if inside[pos] {
				t.Fatalf("位置 %d 被重复占用", pos)
			}
			inside[pos] = true
			stack = append(stack, children[cur]...)
		}
		for index := tree.pos[node]; index < tree.pos[node]+tree.size[node]; index++ {
			if !inside[index] {
				t.Fatalf("子树区间 [%d,%d) 里位置 %d 不属于 %d 的后代", tree.pos[node], tree.pos[node]+tree.size[node], index, node)
			}
		}
	}
}

func TestStar(t *testing.T) {
	count := 8
	values := make([]int64, count)
	edges := [][2]int{}
	for index := 1; index < count; index++ {
		values[index] = int64(index + 1)
		edges = append(edges, [2]int{0, index})
	}
	values[0] = 1
	tree, err := New(count, values, edges)
	if err != nil {
		t.Fatalf("New 出错：%v", err)
	}
	parent, depth, children := buildBrute(count, edges)
	for node := 0; node < count; node++ {
		got, _ := tree.SubtreeSum(node)
		want := bruteSubtreeSum(values, children, node)
		if got != want {
			t.Fatalf("subtree(%d) = %d，暴力是 %d", node, got, want)
		}
		for other := 0; other < count; other++ {
			got, _ := tree.PathSum(node, other)
			want := brutePathSum(values, parent, depth, node, other)
			if got != want {
				t.Fatalf("path(%d,%d) = %d，暴力是 %d", node, other, got, want)
			}
		}
	}
}

func TestRandomAgainstBrute(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	random := rng.Intn
	for round := 0; round < 200; round++ {
		count := 1 + random(60)
		values := make([]int64, count)
		edges := [][2]int{}
		for index := 0; index < count; index++ {
			values[index] = int64(random(100) - 50)
			if index > 0 {
				edges = append(edges, [2]int{random(index), index})
			}
		}
		tree, err := New(count, values, edges)
		if err != nil {
			t.Fatalf("New 出错：%v", err)
		}
		parent, depth, children := buildBrute(count, edges)
		for query := 0; query < 200; query++ {
			first, second := random(count), random(count)
			gotPath, err := tree.PathSum(first, second)
			if err != nil {
				t.Fatalf("PathSum 出错：%v", err)
			}
			wantPath := brutePathSum(values, parent, depth, first, second)
			if gotPath != wantPath {
				t.Fatalf("第 %d 棵树 path(%d,%d) = %d，暴力是 %d", round, first, second, gotPath, wantPath)
			}
			gotSub, err := tree.SubtreeSum(first)
			if err != nil {
				t.Fatalf("SubtreeSum 出错：%v", err)
			}
			wantSub := bruteSubtreeSum(values, children, first)
			if gotSub != wantSub {
				t.Fatalf("第 %d 棵树 subtree(%d) = %d，暴力是 %d", round, first, gotSub, wantSub)
			}
			// 同一输入重复查询结果相同。
			again, _ := tree.PathSum(first, second)
			if again != gotPath {
				t.Fatalf("重复查询结果不一致：%d != %d", again, gotPath)
			}
		}
	}
}

func TestDeepChain(t *testing.T) {
	count := 100000
	values := make([]int64, count)
	for index := range values {
		values[index] = int64(index%7) - 3
	}
	tree, err := New(count, values, chain(count))
	if err != nil {
		t.Fatalf("New 出错：%v", err)
	}
	parent, depth, _ := buildBrute(count, chain(count))
	for query := 0; query < 100000; query++ {
		first := (query * 37) % count
		second := (query * 91) % count
		got, err := tree.PathSum(first, second)
		if err != nil {
			t.Fatalf("PathSum 出错：%v", err)
		}
		if query%997 == 0 {
			want := brutePathSum(values, parent, depth, first, second)
			if got != want {
				t.Fatalf("path(%d,%d) = %d，暴力是 %d", first, second, got, want)
			}
		}
	}
	if sum, _ := tree.SubtreeSum(count / 2); sum != func() int64 {
		total := int64(0)
		for index := count / 2; index < count; index++ {
			total += values[index]
		}
		return total
	}() {
		t.Fatalf("深链子树和不对：%d", sum)
	}
}

func TestQueryBadInput(t *testing.T) {
	tree, err := New(3, []int64{1, 2, 3}, [][2]int{{0, 1}, {0, 2}})
	if err != nil {
		t.Fatalf("New 出错：%v", err)
	}
	if _, err := tree.SubtreeSum(-1); err != ErrInput {
		t.Fatalf("负数节点应该报 ErrInput：%v", err)
	}
	if _, err := tree.SubtreeSum(3); err != ErrInput {
		t.Fatalf("越界节点应该报 ErrInput：%v", err)
	}
	if _, err := tree.PathSum(0, 3); err != ErrInput {
		t.Fatalf("越界节点应该报 ErrInput：%v", err)
	}
	if _, err := tree.PathSum(-1, 0); err != ErrInput {
		t.Fatalf("负数节点应该报 ErrInput：%v", err)
	}
}

func TestMalformedEdges(t *testing.T) {
	if _, err := New(3, []int64{1, 1, 1}, [][2]int{{0, 1}, {2, 1}}); err != ErrInput {
		t.Fatalf("重复父亲应该报 ErrInput：%v", err)
	}
	if _, err := New(3, []int64{1, 1, 1}, [][2]int{{1, 2}}); err != ErrInput {
		t.Fatalf("不连通的边表应该报 ErrInput：%v", err)
	}
	if _, err := New(3, []int64{1, 1, 1}, [][2]int{{0, 0}}); err != ErrInput {
		t.Fatalf("自环应该报 ErrInput：%v", err)
	}
}
