package gee

import (
	"log"
	"net/http"
)

type HandlerFunc func(*context) //归类handler

type Engine struct {
	*RouterGroup
	router *router
	group  []*RouterGroup
}

type RouterGroup struct {
	prefix     string
	middewares []HandlerFunc
	parent     *RouterGroup
	engine     *Engine //将框架中所有的资源统一交给engine进行管理
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.group = []*RouterGroup{engine.RouterGroup}
	return engine

}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
} //添加路由

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.group = append(engine.group, newGroup)
	return newGroup
}
