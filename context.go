package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	Params     map[string]string //动态路由解析之后的键值对
	StatusCode int
	handlers   []HandlerFunc //中间件
	index      int
	engine     *Engine //方便访问engine中的html模板
} //存储所有请求响应信息

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1, //记录运行到第几个中间件
	}
} // 新建context

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
} //返回解析好的from数据中key所代表的项

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
} //返回url查询中key所对于的值

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
} //返回http状态码

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
} //在head中添加键值对

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
} //将请求作为string返回
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
} //将请求作为json解析
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
} //写入data数据
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
} //将数据作为html解析
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
} //取解析后的动态路由对应的值

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	} //递归调用,保证中间件中Next方法后的部分在Handler之后运行
} //顺序执行中间件

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
} //报错
