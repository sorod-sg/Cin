package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node //路由树的根节点
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

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
} //匹配*之前的url

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) handle(c *context) {
	key := c.Method + "-" + c.Path
	if handler, OK := r.handlers[key]; OK {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND :%s\n", c.Path)
	}
}
