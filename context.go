package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
} //上下文

func newContext(w http.ResponseWriter, req *http.Request) *context {
	return &context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

func (c *context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}
func (c *context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}
func (c *context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}
func (c *context) HTML(code int, html string) {
	c.SetHeader("content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
