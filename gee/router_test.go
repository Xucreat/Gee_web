package gee

import (
	"fmt"
	"reflect"
	"testing"
)

/* 简单的 Go 语言测试代码，
用于测试一个基本的路由器（router）功能。 */

// 运行测试：在命令行中运行测试命令，通常是 go test。
// 这将编译并执行测试文件中的所有测试函数，并输出结果。

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

/*
	该函数用于解析路由路径，返回路由的模式和参数。测试了三种情况：

"/p/:name" 应该返回 []string{"p", ":name"}
"/p/*" 应该返回 []string{"p", "*"}
"/p/*name/*" 应该返回 []string{"p", "*name"}
*/
func TestParsePattern(t *testing.T) {
	// 使用了 reflect.DeepEqual 函数来比较切片是否相等，如果不相等则表示测试失败。
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

/*
	该函数用于匹配请求的路径并返回匹配的路由节点和参数。

通过调用 r.getRoute("GET", "/hello/geektutu") 测试是否能够正确匹配路径 "/hello/geektutu"。
如果匹配成功，应该返回路由节点和参数，其中路由节点的模式应该是 "/hello/:name"，参数中 name 应该是 "geektutu"。
*/
func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/xhl")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if ps["name"] != "xhl" {
		t.Fatal("name should be equal to 'xhl'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])

}
