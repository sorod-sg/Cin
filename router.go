package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       //路由树的根节点
	handlers map[string]HandlerFunc //路由的处理handler
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

//匹配*之前的part,将*之前的part切片
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != " " {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{} //若root节点不存在,创建新的root节点
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.path)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index] //将通配符":"解析为对应路径
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/") //将通配符"*"解析为对应路径
				break
			}
		}

	}
	return n, params //返回通配符的匹配
}
func (r *router) handle(c *context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.path
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND :%s\n", c.Path)
	}
}
