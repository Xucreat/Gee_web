package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	// 将路径模式分解为部分，并返回一个字符串切片 parts，每个元素代表路径中的一个部分。
	parts := parsePattern(pattern)
	// 键通常由 HTTP 方法和路径模式组成
	key := method + "-" + pattern

	// 不存在指定 HTTP 方法的根节点，则创建一个新的根节点，并将其存储在 r.roots 中，以 HTTP 方法作为键
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	// 将路径模式的部分插入到路由树中。这将在路由树中创建或更新与路径模式匹配的节点。
	r.roots[method].insert(pattern, parts, 0)

	// 将处理函数 handler 存储在 r.handlers 中，以路径模式和 HTTP 方法组合成的键 key 作为索引。
	// 建立了路径模式和处理函数之间的映射关系
	r.handlers[key] = handler
}

// func (r *router) handle(c *Context) {
// 	key := c.Method + "-" + c.Path
// 	if handler, ok := r.handlers[key]; ok {
// 		handler(c)
// 	} else {
// 		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
// 	}
// }

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) // 将给的路径模式分解为部分
	params := make(map[string]string)
	root, ok := r.roots[method] // 某个HTTP方法的前缀树的根节点

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0) // 前缀树中与给的路径匹配上的路径

	if n != nil {
		// 前缀树中与给的路径匹配上的路径模式分解为部分，
		// 并返回一个字符串切片 parts，每个元素代表路径中的一个部分
		parts := parsePattern(n.pattern)

		// 解析了:和*两种匹配符的参数，返回一个 map
		for index, part := range parts {
			// 如/p/go/doc匹配到/p/:lang/doc，解析结果为：{lang: "go"}
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			// /static/css/geektutu.css匹配到/static/*filepath，
			// 解析结果为{filepath: "css/geektutu.css"}。
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		// 在调用匹配到的handler前，将解析出来的路由参数赋值给了c.Params
		// 能够在handler中，通过Context对象访问到具体的值。
		r.handlers[key](c)

	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
