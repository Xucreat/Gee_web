package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// implements the interface of ServeHTTP -> as middleware?
type Engine struct {
	*RouterGroup
	// router map[string]HandlerFunc
	router *router
	groups []*RouterGroup // store all groups
	// 将所有的模板加载进内存
	htmlTemplates *template.Template // for html render
	// 所有的自定义模板渲染函数
	funcMap template.FuncMap // for html render
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

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
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
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

// 设置自定义渲染函数funcMap的方法
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// 加载模板的方法
func (engine *Engine) LoadHTMLGlob(pattern string) {
	// 调用 template.New("") 创建一个新的命名模板，并且给它一个空字符串作为名称。这个命名模板将用于解析和存储匹配模式的所有模板。
	// 调用 .Funcs(engine.funcMap) 方法将 Engine 结构体中的 funcMap 添加到模板中。funcMap 是一个包含自定义函数的映射（map），这些函数可以在模板中使用。
	// 调用 .ParseGlob(pattern) 方法，使用给定的模式去匹配并解析所有符合条件的模板文件。所有匹配到的模板文件都会被解析并添加到之前创建的命名模板中。
	// 将上述操作放在 template.Must 函数中，确保在解析模板时如果发生任何错误，程序会立即崩溃并报告错误信息。这是 Go 语言中一种常见的错误处理方式，尤其在初始化阶段使用。
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

/*静态文件(Serve Static Files)*/
// create static handler
func (group *RouterGroup) createStasticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	// 使用 http.StripPrefix 去除请求路径中的 absolutePath 部分，
	// 然后创建一个文件服务器 http.FileServer 来处理静态文件请求。
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) { // 返回一个匿名函数，作为处理静态文件请求的处理函数。
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		// 文件存在，调用 fileServer.ServeHTTP 方法处理文件请求，将文件内容返回给客户端。
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// serve static files
// Static方法是暴露给用户的
// 用户可以将磁盘上的某个文件夹root映射到路由relativePath。
// 例如:用户访问localhost:9999/assets/js/xhl.js，最终返回/usr/xhl/blog/static/js/xhl.js。
func (group *RouterGroup) Static(relativePath string, root string) {
	//  http.Dir(root) 将文件系统目录 root 转换为 http.FileSystem 接口。
	/*
			type FileSystem interface {
		    Open(name string) (File, error)
			}
	*/
	handler := group.createStasticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath") // relativePath 是 /assets, urlPattern == "/assets/*filepath"
	// 若不写上两行代码,fileHandler.ServeHTTP 会把req.url.path 作为文件路径

	// Register GET handlers
	group.GET(urlPattern, handler)
}
