// 命令行入口：heavy-light-tree --values 点权 --edges 边表 [--path u,v] [--subtree 根]
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	heavylight "example.com/heavylight"
)

func parseValues(text string) ([]int64, error) {
	if text == "" {
		return nil, fmt.Errorf("values 为空")
	}
	parts := strings.Split(text, ",")
	values := make([]int64, 0, len(parts))
	for _, part := range parts {
		value, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("非法点权 %q: %w", part, err)
		}
		values = append(values, value)
	}
	return values, nil
}

func parseEdges(text string) ([][2]int, error) {
	if text == "" {
		return nil, nil
	}
	parts := strings.Split(text, ",")
	edges := make([][2]int, 0, len(parts))
	for _, part := range parts {
		pair := strings.Split(strings.TrimSpace(part), "-")
		if len(pair) != 2 {
			return nil, fmt.Errorf("非法边 %q", part)
		}
		parent, err := strconv.Atoi(strings.TrimSpace(pair[0]))
		if err != nil {
			return nil, fmt.Errorf("非法边 %q: %w", part, err)
		}
		child, err := strconv.Atoi(strings.TrimSpace(pair[1]))
		if err != nil {
			return nil, fmt.Errorf("非法边 %q: %w", part, err)
		}
		edges = append(edges, [2]int{parent, child})
	}
	return edges, nil
}

func parsePair(text string) (int, int, error) {
	pair := strings.Split(strings.TrimSpace(text), ",")
	if len(pair) != 2 {
		return 0, 0, fmt.Errorf("非法路径 %q", text)
	}
	first, err := strconv.Atoi(strings.TrimSpace(pair[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("非法路径 %q: %w", text, err)
	}
	second, err := strconv.Atoi(strings.TrimSpace(pair[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("非法路径 %q: %w", text, err)
	}
	return first, second, nil
}

func run() error {
	valuesFlag := flag.String("values", "", "点权，逗号分隔")
	edgesFlag := flag.String("edges", "", "边表，逗号分隔的 parent-child")
	pathFlag := flag.String("path", "", "路径求和 u,v")
	subtreeFlag := flag.Int("subtree", -1, "子树求和的根")
	flag.Parse()

	values, err := parseValues(*valuesFlag)
	if err != nil {
		return err
	}
	edges, err := parseEdges(*edgesFlag)
	if err != nil {
		return err
	}
	tree, err := heavylight.New(len(values), values, edges)
	if err != nil {
		return err
	}
	if *pathFlag != "" {
		first, second, err := parsePair(*pathFlag)
		if err != nil {
			return err
		}
		sum, err := tree.PathSum(first, second)
		if err != nil {
			return err
		}
		fmt.Printf("path(%d,%d)=%d\n", first, second, sum)
	}
	if *subtreeFlag >= 0 {
		sum, err := tree.SubtreeSum(*subtreeFlag)
		if err != nil {
			return err
		}
		fmt.Printf("subtree(%d)=%d\n", *subtreeFlag, sum)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
