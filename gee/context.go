package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Method     string
	Path       string
	Writer     http.ResponseWriter
	Request    *http.Request
	Params     map[string]string
	StatusCode int
	//用户请求要执行的函数链
	handlers []HandleFunc
	index    int
}

func NewContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		Method:  request.Method,
		Path:    request.URL.Path,
		Writer:  writer,
		Request: request,
		index:   -1,
	}
}
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}
func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// 没用过先不写
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 没用过先不写
func (c *Context) String(code int, value string) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(value)))
}

func (c *Context) JSON(httpCode int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(httpCode)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}
func (c *Context) HTML(httpCode int, name string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(httpCode)
	c.Writer.Write([]byte(name))
}
func (c *Context) Fail(httpCode int, message string) {
	c.String(httpCode, message)
}
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}
