package gee

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
)

// 提供一些简单的中间件

// Logger 提供一个日志中间件
func Logger() HandlerFunc {
	return func(c *Context) {
		tim := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(tim))
	}
}

// trace 提供堆栈跟踪信息
func trace(message string) string {
	var res strings.Builder
	var pc []uintptr
	n := runtime.Callers(3, pc)
	res.WriteString(message)
	res.WriteString("\nTraceback:\n")
	for _, v := range pc[:n] {
		fun := runtime.FuncForPC(v)
		file, line := fun.FileLine(v)
		res.WriteString(fmt.Sprintf("%s:%d %s\n", file, line, fun.Name()))
	}
	res.WriteString("\n\n")
	return res.String()
}

// Recovery 提供一个错误处理中间件
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := err.(string)
				log.Printf("Recovery from panic: %s\n", message)
				log.Printf("%s\n", trace(message))
				c.Fail(500, message)
			}
		}()
		c.Next()
	}
}
