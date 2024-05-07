package gee

import (
	"net/http"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// implements the interface of ServeHTTP -> as middleware?
type Engine struct {
	// router map[string]HandlerFunc
	router *router
}

// contructor of gee.Engine
func New() *Engine {
	// return &Engine{router: make(map[string]HandlerFunc)}
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	// key := method + "-" + pattern
	// engine.router[key] = handler
	engine.router.addRoute(method, pattern, handler)
}

// defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// defines the method to start a http request
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
} // package the http.LinstenAndServer method

// Engine implements http.Handler (method ServeHTTP)
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// key := req.Method + "-" + req.URL.Path
	// if handler, ok := engine.router[key]; ok {
	// 	handler(w, req)
	// } else {
	// 	fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL) // 将 404 NOT FOUND 的消息写入 HTTP 响应，消息中包含了请求的 URL。
	// }
	c := newContext(w, req)
	engine.router.handle(c)
}
