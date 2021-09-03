package gee

import "strings"

type node struct {
	pattern  string //待匹配路由
	part     string //路由中的一部分
	chlidren []*node
	isWild   bool //是否精准匹配
} //前缀树路由的节点

func (n *node) matchChlid(part string) *node {
	for _, child := range n.chlidren {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
} //寻找匹配的下一子节点

func (n *node) matchChlidren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.chlidren {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
} //所有匹配成功的节点

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChlid(part)
	if child == nil {
		child = &node{part: part, isWild: (part[0] == ':' || part[0] == '*')}
	}
	child.insert(pattern, parts, height+1)
} //插入节点

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	chlidren := n.matchChlidren(part)

	for _, child := range chlidren {
		rusult := child.search(parts, height+1)
		if rusult != nil {
			return rusult
		}
	}
	return nil
}
