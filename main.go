package main

import (
	"net/http"

	"gee"
)

func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		// fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/hello", func(c *gee.Context) {
		// for k, v := range req.Header {
		// 	fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		// }

		// expect /hello?name=xhl
		c.String(http.StatusOK, "hello %s, you're at %s", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9090")
}
