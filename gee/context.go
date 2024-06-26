package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/* 封装 Request 和 Response，提供JSON、HTML等返回类型的支持 */
/* 设计Context结构，扩展性和复杂性留在内部，对外简化接口。 */
type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	// 将解析后的参数存储到Params中，
	// 通过c.Param("lang")的方式获取到对应的值。
	Params map[string]string

	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int

	engine *Engine
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}

func (c *Context) Param(key string) string {
	value := c.Params[key]
	return value
}

// ?
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// 接收到请求后，应查找所有应作用于该路由的中间件，保存在Context中，依次进行调用。
// index是记录当前执行到第几个中间件，当在中间件中调用Next方法时，控制权交给了下一个中间件，
// 直到调用到最后一个中间件，
// 然后再从后往前，调用每个中间件在Next方法之后定义的部分。
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}

}
