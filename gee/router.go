package gee

import (
	"net/http"
	"strings"
)

// 将路由相关的方法提取出来

type Router struct {
	handlers map[string]HandlerFunc
	roots map[string]*node
}

func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]HandlerFunc),
		roots:    make(map[string]*node),
	}
}

// 用于解析url
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	res := make([]string, 0)

	for _, v := range vs {
		if v != "" {
			res = append(res, v)
			if v[0] == '*' {
				break
			}
		}
	}

	return res
}

func (r *Router) addRouter(method, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler

	// 在roots中注册该url的逻辑

	parts := parsePattern(pattern)

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
}

func (r *Router) handle(c *Context) {
	n, params := r.getRouter(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		http.Error(c.Writer, "404 NOT FOUND\n", 404)
	}

	c.Next()
}

// 获取满足某个path信息的节点
func (r *Router) getRouter(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)

	root, ok := r.roots[method]
	if !ok { // 这个请求列表完全为空
		return nil, nil
	}

	n := root.search(searchParts, 0)
	if n == nil {
		return nil, nil
	}

	parts := parsePattern(n.pattern)

	for k, v := range parts {
		if v[0] == ':' {
			params[v[1:]] = searchParts[k]
		}
		if v[0] == '*' && len(v) > 1 {
			params[v[1:]] = strings.Join(searchParts[k:], "/")
			break
		}
	}

	return n, params
}
