package gee

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
	}
}
