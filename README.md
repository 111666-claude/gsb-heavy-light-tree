# heavy-light-tree

树链剖分（路径求和 + 子树求和），只用 Go 标准库。

```
go test ./...
go run ./cmd/heavy-light-tree --values 1,2,3,4 --edges '0-1,1-2,1-3' --path 2,3 --subtree 1
go run ./cmd/heavy-light-tree --values 1,2,3,4 --edges '0-1,1-2,1-3' --subtree 2
```

## 口径（README 为准）

- **位置必须是「重儿子优先」的 DFS 序**：这样每条链连续、每棵子树也连续（子树区间是 `[pos, pos+size)`）。
  现在的位置按链编号顺序拼接，子树区间不连续：`--edges '0-1,1-2,1-3' --values 1,2,3,4 --subtree 1`
  按口径是 `subtree(1)=9`，现在给 `subtree(1)=2`。
- **重儿子要选子树最大的孩子**：`measure` 里按 `size` 比较（现在写着比较但位置拼接没用上，
  必须让 `heavy`、`head` 与位置分配三者一致）。
- 路径查询：沿两条链往上跳，链内用前缀和/线段树求区间和，最后加上 LCA 那一段；结果要与逐点向上走的暴力实现一致。
- 不变量：`SubtreeSum(u)` 等于以 u 为根的子树所有点权之和；`PathSum(u,v)` 等于 u 到 v 简单路径上的点权和；
  同一输入重复查询结果相同；非法节点报 `ErrInput`。
- 规模：10 万个点的深链（每个点只有一个孩子）要在 1 秒内建好并回答 10 万次查询，
  内存 O(n)；不许用递归深度超过 n 的实现（深链必须走迭代或保证不爆栈）。

## 现在的行为

- `New` 里位置 `order[index] = index`（按编号顺序），不是重儿子优先的 DFS 序，子树区间因此不连续。
- `children` 每次遍历全表找 `parent == -1` 的节点，是占位实现；`heavy` 选出来后也没参与位置分配。
- `PathSum` 用逐点向上走的朴素方式，结果是对的，但复杂度 O(深度)。

## 计划新增的能力

```
--values 点权     逗号分隔
--edges 边表      逗号分隔的 parent-child
--path u,v        路径求和
--subtree 根      子树求和
```

## 输出

```
path(2,3)=7
subtree(1)=9
```

## 目录

```
hld.go                  建树与查询
cmd/heavy-light-tree    命令行入口
hld_test.go             go test 用例
```
