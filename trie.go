package gee

import (
	"runtime"

	"honnef.co/go/tools/pattern"
)

type node struct {
	pattern  string //待匹配路由
	part     string
	chlidren []*node
	isWild   bool
}

func (n *node) matchChlid(part string) *node {
	for _, child := range n.chlidren {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChlidren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.chlidren {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

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
}
