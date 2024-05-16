package gee

import (
	"log"
	"net/http"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// implements the interface of ServeHTTP -> as middleware?
type Engine struct {
	*RouterGroup
	// router map[string]HandlerFunc
	router *router
	groups []*RouterGroup // store all groups
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine instance
}

// contructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
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
