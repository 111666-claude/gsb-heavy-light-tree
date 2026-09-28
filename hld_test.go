package heavylight

import "testing"

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
