package gee

// 针对使用场景，封装Req和Writer，简化使用

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type H map[string]interface{}

var mux sync.Mutex

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Method string
	Path   string
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	idx      int
	// 让Context也可以访问到engine
	engine *Engine
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:   w,
		Req:      r,
		Method:   r.Method,
		Path:     r.URL.Path,
		Params:   map[string]string{},
		handlers: []HandlerFunc{},
		idx:      -1,
	}
}

func (c *Context) Next() {
	c.idx++
	for i := 0; i < len(c.handlers); i++ {
		c.handlers[i](c)
	}
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Status(code int) {
	if c.StatusCode == code { // 避免重复设置状态码
		return
	}
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, value interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(value); err != nil { //将value以json解码，如果成功则err为nil
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, htmlFile string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if er := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, htmlFile, data); er != nil {
		c.Fail(500, er.Error())
	}
}
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// Fail 返回错误信息
func (c *Context) Fail(i int, s string) {
	c.Writer.WriteHeader(i)
	c.Writer.Write([]byte(s))
}
