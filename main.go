package main

import (
	"gee"
	"net/http"
)

/*
$ curl "http://localhost:9999"
Hello Xhl
$ curl "http://localhost:9999/panic"
{"message":"Internal Server Error"}
$ curl "http://localhost:9999"
Hello Xhl

>>> log
2024/05/30 12:29:49 [200] / in 16.2µs
2024/05/30 12:30:05 runtime error: index out of range [100] with length 1
Traceback:
        /opt/go/src/runtime/panic.go:914
        /opt/go/src/runtime/panic.go:114
        /home/bx/code_files/Gee_web/main.go:16
        /home/bx/code_files/Gee_web/gee/context.go:106
        /home/bx/code_files/Gee_web/gee/recovery.go:37
        /home/bx/code_files/Gee_web/gee/context.go:106
        /home/bx/code_files/Gee_web/gee/logger.go:15
        /home/bx/code_files/Gee_web/gee/context.go:106
        /home/bx/code_files/Gee_web/gee/router.go:118
        /home/bx/code_files/Gee_web/gee/gee.go:100
        /opt/go/src/net/http/server.go:2939
        /opt/go/src/net/http/server.go:2010
        /opt/go/src/runtime/asm_amd64.s:1651

2024/05/30 12:30:05 [500] /panic in 1.419724ms
*/

func main() {
	r := gee.Default()
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello Xhl\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"xhlxhl"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
