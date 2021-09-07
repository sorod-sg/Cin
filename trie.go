package gee

import "strings"

//前缀树路由的节点
type node struct {
	part     string //待匹配路由
	path     string //路由路径
	chlidren []*node
	isWild   bool //是否为通配符节点
}

//寻找匹配的下一子节点
func (n *node) matchChlid(part string) *node {
	for _, child := range n.chlidren {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

//所有匹配成功的节点
func (n *node) matchChlidren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.chlidren {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//插入节点 从height开始进行dfs 从root节点开始插入的话输入height为0
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.path = pattern
		return
	}
	part := parts[height]
	child := n.matchChlid(part)
	if child == nil {
		child = &node{part: part, isWild: (part[0] == ':' || part[0] == '*')}
	}
	child.insert(pattern, parts, height+1)
}

//从parts中的第height节点开始搜索是否有part路径,默认搜索height为0
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.path == "" {
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
