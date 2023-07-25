package gee

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string][]HandleFunc
}

func newRouter() *router {
	return &router{make(map[string]*node), make(map[string][]HandleFunc)}
}

// 解析路径为切片，“/”分割；遇到“*”如 /img/*p/a=> {img,*p,a}，
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0, len(vs))
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
func foo(key string, f HandleFunc) {
	fmt.Printf("%s Function name: %s\n", key, getFunctionName(f))
}

func getFunctionName(f HandleFunc) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// 添加路由
func (r *router) addRoute(method string, p string, handleFunc []HandleFunc) {
	parts := parsePattern(p) //解析路由
	key := method + "-" + p
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(p, parts, 0)
	if _, ok1 := r.handlers[key]; ok1 {
		panic("已注册相同路由" + key)
	}
	r.handlers[key] = handleFunc
}

// 获取路由
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil { //搜索节点不为空
		parts := parsePattern(n.pattern) //解析路由
		for i, part := range parts {     //遍历每一层路由
			if part[0] == ':' { //动态路由，添加到动态路径参数的映射。如/:name匹配到/1，m[name]=1
				params[part[1:]] = searchParts[i]
			}
			if part[0] == '*' && len(part) > 1 { //动态路径如/c/*a匹配到/c/a/b，m[a]=/a/b
				params[part[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}
func (r *router) handle(ctx *Context) {
	n, params := r.getRoute(ctx.Method, ctx.Path)
	if n != nil {
		ctx.Params = params
		key := ctx.Method + "-" + n.pattern
		ctx.handlers = r.handlers[key]
	} else {
		ctx.handlers = append(ctx.handlers, func(c *Context) {
			ctx.String(http.StatusNotFound, "404 Not Found Router:"+ctx.Path)
		})
	}
	ctx.Next()
}
