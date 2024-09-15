package main

import (
	"gee"
	"net/http"
)

// 测试

func main() {
	r := gee.Default()
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello Jesse\n")
	})
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"Jesse"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
