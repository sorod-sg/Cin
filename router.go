package gee

import (
	"log"
	"net/http"
)

type router struct {
	handlers map[string] HandlerFunc
}

func newRouter() *router{
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRoute (method string ,pattern string ,handler HandlerFunc){
	log.Printf("Route %4s - %s" , method ,pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}
func (r *router)handle (c *context){
	key := c.Method + "-" + c.Path
	if handler ,OK :=r.handlers[key] ;OK{
		handler(c)
	}else{
		c.String(http.StatusNotFound ,"404 NOT FOUND :%s\n" ,c.Path)
	}
}