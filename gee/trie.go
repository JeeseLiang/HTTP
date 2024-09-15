package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  // 是否完整路由(isEnd),是则为完整url,否则为空
	part     string  // 该节点对应的路由
	children []*node // 子节点
	isWild   bool    // 是否精确匹配
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// 找到某层第一个匹配的节点 用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.isWild || child.part == part {
			return child
		}
	}
	return nil
}

// 这个函数跟matchChild有点像，但它是返回所有匹配的子节点，原因是它的场景是用以查找
// 它必须返回所有可能的子节点来进行遍历查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0) // 存储该部分匹配成功的路由集合
	for _, child := range n.children {
		if child.isWild || child.part == part { // 该部分匹配成功
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 插入一个路由
func (n *node) insert(pattern string, parts []string, iter int) {
	if iter == len(parts) { // 匹配完成
		n.pattern = pattern
		return
	}

	child := n.matchChild(parts[iter])

	if child == nil {
		child = &node{
			part:   parts[iter],
			isWild: parts[iter][0] == ':' || parts[iter][0] == '*',
		}
		n.children = append(n.children, child)
	}

	child.insert(pattern, parts, iter+1)
}

func (n *node) search(parts []string, iter int) *node {
	if len(parts) == iter || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	children := n.matchChildren(parts[iter])
	// 获得所有可能的路径

	for _, v := range children {
		res := v.search(parts, iter+1)
		if res != nil {
			return res
			// 找到符合的可以直接返回
		}
	}

	return nil
}
