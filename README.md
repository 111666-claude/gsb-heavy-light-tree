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

## 实现

- `New` 用迭代 DFS 求子树 size 并选重儿子，再按重儿子优先的迭代 DFS 序分配位置，
  深链不递归、不爆栈，内存 O(n)。
- `heavy`、`head` 与位置分配三者一致：沿重链连续编号，轻儿子入栈稍后处理，
  因此每条链连续、每棵子树是区间 `[pos, pos+size)`。
- `PathSum` 沿两条链向上跳，链内用前缀和求区间和，最后补上 LCA 那一段，复杂度 O(log² n)。

## 计划新增的能力

```
--values 点权     逗号分隔
--edges 边表      逗号分隔的 parent-child
--path u,v        路径求和
--subtree 根      子树求和
```

## 输出

```
path(2,3)=9
subtree(1)=9
```

## 目录

```
hld.go                  建树与查询
cmd/heavy-light-tree    命令行入口
hld_test.go             go test 用例
```
