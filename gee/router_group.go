package gee

type RouterGroup struct {
	prefix      string
	middlewares []HandleFunc // support middleware
	parent      *RouterGroup // support nesting
	engine      *Engine      // all groups share a Engine instance
}

func (g *RouterGroup) Use(handleFunc ...HandleFunc) {
	g.middlewares = append(g.middlewares, handleFunc...)
}
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		middlewares: group.middlewares,
		prefix:      group.prefix + prefix,
		parent:      group,
		engine:      engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (g *RouterGroup) addRoute(method string, r string, handleFunc HandleFunc) {
	// 将当前路由分组中定义的中间件函数列表复制一份
	middlewares := make([]HandleFunc, len(g.middlewares))
	copy(middlewares, g.middlewares)

	// 将路由处理函数添加到处理函数列表中

	middlewares = append(middlewares, handleFunc)
	pattern := g.prefix + r
	//打印函数
	key := method + "-" + pattern
	foo(key, handleFunc)
	g.engine.router.addRoute(method, pattern, middlewares)
}
func (g *RouterGroup) GET(r string, handleFunc HandleFunc) {
	g.addRoute("GET", r, handleFunc)
}
func (g *RouterGroup) PUT(r string, handleFunc HandleFunc) {
	g.addRoute("PUT", r, handleFunc)
}
func (g *RouterGroup) POST(r string, handleFunc HandleFunc) {
	g.addRoute("POST", r, handleFunc)
}
func (g *RouterGroup) DELETE(r string, handleFunc HandleFunc) {
	g.addRoute("DELETE", r, handleFunc)
}
