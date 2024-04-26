package gee

import "net/http"

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
	// response info
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

func (c.Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}
