package heavylight

import (
	"math/rand"
	"testing"
	"time"
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

func bruteSubtree(values []int64, children [][]int, root int) int64 {
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

func brutePath(values []int64, parent, depth []int, first, second int) int64 {
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

func buildBrute(count int, edges [][2]int) ([][]int, []int, []int) {
	children := make([][]int, count)
	parent := make([]int, count)
	depth := make([]int, count)
	for index := range parent {
		parent[index] = -1
	}
	for _, edge := range edges {
		children[edge[0]] = append(children[edge[0]], edge[1])
		parent[edge[1]] = edge[0]
	}
	queue := []int{0}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, child := range children[node] {
			depth[child] = depth[node] + 1
			queue = append(queue, child)
		}
	}
	return children, parent, depth
}

func TestHeavyChoice(t *testing.T) {
	// 0 的孩子：1（子树大小 2）与 2（子树大小 1），重儿子必须是 1。
	tree, err := New(4, []int64{5, 6, 7, 8}, [][2]int{{0, 1}, {0, 2}, {1, 3}})
	if err != nil {
		t.Fatalf("New 出错：%v", err)
	}
	if tree.heavy[0] != 1 {
		t.Fatalf("0 的重儿子应该是 1：%d", tree.heavy[0])
	}
	if tree.heavy[1] != 3 {
		t.Fatalf("1 的重儿子应该是 3：%d", tree.heavy[1])
	}
	if tree.heavy[2] != -1 || tree.heavy[3] != -1 {
		t.Fatalf("叶子没有重儿子：%d %d", tree.heavy[2], tree.heavy[3])
	}
	// heavy、head 与位置三者一致：重儿子与父同链且位置相邻。
	for node := 0; node < tree.count; node++ {
		if child := tree.heavy[node]; child >= 0 {
			if tree.head[child] != tree.head[node] {
				t.Fatalf("重儿子 %d 应与 %d 同链", child, node)
			}
			if tree.pos[child] != tree.pos[node]+1 {
				t.Fatalf("重儿子 %d 应紧跟 %d 的位置", child, node)
			}
		} else if tree.head[node] != node && tree.head[tree.head[node]] != tree.head[node] {
			t.Fatalf("链头标记不一致：%d", node)
		}
	}
}

func TestSubtreeInterval(t *testing.T) {
	// 题目用例：subtree(1) 只含 {1,3}，应为 14。
	tree, err := New(4, []int64{5, 6, 7, 8}, [][2]int{{0, 1}, {0, 2}, {1, 3}})
	if err != nil {
		t.Fatalf("New 出错：%v", err)
	}
	if sum, _ := tree.SubtreeSum(1); sum != 14 {
		t.Fatalf("subtree(1) 应该是 14：%d", sum)
	}
	// 子树区间 [pos, pos+size) 恰好覆盖所有后代。
	children, _, _ := buildBrute(4, [][2]int{{0, 1}, {0, 2}, {1, 3}})
	for root := 0; root < 4; root++ {
		inside := make([]bool, 4)
		stack := []int{root}
		for len(stack) > 0 {
			node := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			inside[node] = true
			stack = append(stack, children[node]...)
		}
		for node := 0; node < 4; node++ {
			inInterval := tree.pos[node] >= tree.pos[root] && tree.pos[node] < tree.pos[root]+tree.size[root]
			if inInterval != inside[node] {
				t.Fatalf("节点 %d 相对子树 %d 的区间归属错误", node, root)
			}
		}
	}
}

func TestStarTree(t *testing.T) {
	count := 6
	values := []int64{1, 2, 3, 4, 5, 6}
	edges := [][2]int{}
	for index := 1; index < count; index++ {
		edges = append(edges, [2]int{0, index})
	}
	tree, err := New(count, values, edges)
	if err != nil {
		t.Fatalf("New 出错：%v", err)
	}
	if sum, _ := tree.SubtreeSum(0); sum != 21 {
		t.Fatalf("星形树根子树和应该是 21：%d", sum)
	}
	if sum, _ := tree.SubtreeSum(3); sum != 4 {
		t.Fatalf("叶子子树和应该是 4：%d", sum)
	}
	if sum, _ := tree.PathSum(2, 4); sum != 9 {
		t.Fatalf("path(2,4)=3+1+5 应该是 9：%d", sum)
	}
}

func TestInvalidNode(t *testing.T) {
	tree, err := New(3, []int64{1, 1, 1}, [][2]int{{0, 1}, {1, 2}})
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
		t.Fatalf("越界路径应该报 ErrInput：%v", err)
	}
	if _, err := tree.PathSum(-1, 0); err != ErrInput {
		t.Fatalf("负数路径应该报 ErrInput：%v", err)
	}
}

func TestRandomAgainstBrute(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for trial := 0; trial < 200; trial++ {
		count := 1 + rng.Intn(50)
		values := make([]int64, count)
		for index := range values {
			values[index] = int64(rng.Intn(201) - 100)
		}
		edges := make([][2]int, 0, count-1)
		for node := 1; node < count; node++ {
			edges = append(edges, [2]int{rng.Intn(node), node})
		}
		tree, err := New(count, values, edges)
		if err != nil {
			t.Fatalf("New 出错：%v", err)
		}
		children, parent, depth := buildBrute(count, edges)
		for query := 0; query < 200; query++ {
			first, second := rng.Intn(count), rng.Intn(count)
			pathGot, err := tree.PathSum(first, second)
			if err != nil {
				t.Fatalf("PathSum 出错：%v", err)
			}
			want := brutePath(values, parent, depth, first, second)
			if pathGot != want {
				t.Fatalf("树 %d path(%d,%d)：得到 %d，暴力 %d", trial, first, second, pathGot, want)
			}
			got, err := tree.SubtreeSum(first)
			if err != nil {
				t.Fatalf("SubtreeSum 出错：%v", err)
			}
			want = bruteSubtree(values, children, first)
			if got != want {
				t.Fatalf("树 %d subtree(%d)：得到 %d，暴力 %d", trial, first, got, want)
			}
			// 同一输入重复查询结果相同。
			again, _ := tree.PathSum(first, second)
			againSub, _ := tree.SubtreeSum(first)
			if again != pathGot || againSub != got {
				t.Fatalf("树 %d 重复查询结果不一致", trial)
			}
		}
	}
}

func TestDeepChain(t *testing.T) {
	const count = 100000
	values := make([]int64, count)
	for index := range values {
		values[index] = 1
	}
	start := time.Now()
	tree, err := New(count, values, chain(count))
	if err != nil {
		t.Fatalf("New 出错：%v", err)
	}
	for query := 0; query < 100000; query++ {
		first := query % count
		second := count - 1 - first
		distance := first - second
		if distance < 0 {
			distance = -distance
		}
		if sum, err := tree.PathSum(first, second); err != nil || sum != int64(distance+1) {
			t.Fatalf("path(%d,%d) 错误：%d %v", first, second, sum, err)
		}
		if sum, err := tree.SubtreeSum(first); err != nil || sum != int64(count-first) {
			t.Fatalf("subtree(%d) 错误：%d %v", first, sum, err)
		}
	}
	t.Logf("10 万深链建树 + 20 万次查询耗时 %v", time.Since(start))
}
