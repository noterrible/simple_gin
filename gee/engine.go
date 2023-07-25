package gee

import (
	"net/http"
)

type H map[string]interface{}
type HandleFunc func(c *Context)
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

func Default() *Engine {
	e := New()
	e.Use(Recovery(), Logger())

	return e
}
func New() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
	e.router.handle(c)
}

func (e *Engine) Run(ipAndPort string) error {
	return http.ListenAndServe(ipAndPort, e)
}
