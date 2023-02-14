package route

import (
	"github.com/labstack/echo/v4"
)

// RouterData 路由结构
type RouterData struct {
	title      string
	prefix     string
	permission bool
	data       []*RouterItem
	group      []*RouterData
	router     *echo.Group
}

// RouterItem 路由结构
type RouterItem struct {
	title  string
	method string
	path   string
	name   string
}

// New 新建资源路由
func New(router *echo.Group) *RouterData {
	return &RouterData{
		router: router,
	}
}

// Group 路由分组
func (t *RouterData) Group(prefix string, title string, middle ...echo.MiddlewareFunc) *RouterData {
	group := &RouterData{
		title:  title,
		prefix: prefix,
		router: t.router.Group(prefix, middle...),
	}
	t.group = append(t.group, group)
	return group
}

// Permission 设置权限路由
func (t *RouterData) Permission() *RouterData {
	t.permission = true
	return t
}

// Router 返回原始路由
func (t *RouterData) Router() *echo.Group {
	return t.router
}

// Get 路由
func (t *RouterData) Get(path string, handler echo.HandlerFunc, title string, name string) *echo.Route {
	return t.Add(echo.GET, path, handler, title, name)
}

// Head 路由
func (t *RouterData) Head(path string, handler echo.HandlerFunc, title string, name string) *echo.Route {
	return t.Add(echo.HEAD, path, handler, title, name)
}

// Post 路由
func (t *RouterData) Post(path string, handler echo.HandlerFunc, title string, name string) *echo.Route {
	return t.Add(echo.POST, path, handler, title, name)
}

// Put 路由
func (t *RouterData) Put(path string, handler echo.HandlerFunc, title string, name string) *echo.Route {
	return t.Add(echo.PUT, path, handler, title, name)
}

// Delete 路由
func (t *RouterData) Delete(path string, handler echo.HandlerFunc, title string, name string) *echo.Route {
	return t.Add(echo.DELETE, path, handler, title, name)
}

// Connect 路由
func (t *RouterData) Connect(path string, handler echo.HandlerFunc, title string, name string) *echo.Route {
	return t.Add(echo.CONNECT, path, handler, title, name)
}

// Options 路由
func (t *RouterData) Options(path string, handler echo.HandlerFunc, title string, name string) *echo.Route {
	return t.Add(echo.OPTIONS, path, handler, title, name)
}

// Trace 路由
func (t *RouterData) Trace(path string, handler echo.HandlerFunc, title string, name string) *echo.Route {
	return t.Add(echo.TRACE, path, handler, title, name)
}

// Patch 路由
func (t *RouterData) Patch(path string, handler echo.HandlerFunc, title string, name string) *echo.Route {
	return t.Add(echo.PATCH, path, handler, title, name)
}

// Add 添加路由资源
func (t *RouterData) Add(method string, path string, handler echo.HandlerFunc, title string, name string) *echo.Route {
	item := RouterItem{
		title:  title,
		method: method,
		path:   path,
		name:   name,
	}
	t.data = append(t.data, &item)
	r := t.router.Add(method, path, handler)
	r.Name = item.name
	return r
}

// ParseTree 解析路由为树形
func (t *RouterData) ParseTree(prefix string) any {
	var all []any
	for _, datum := range t.data {
		all = append(all, map[string]any{
			"title":  datum.title,
			"name":   datum.name,
			"method": datum.method,
			"path":   prefix + datum.path,
		})
	}
	for _, item := range t.group {
		gpath := prefix + item.prefix
		all = append(all, item.ParseTree(gpath))
	}
	if t.title == "" {
		return all
	}
	return map[string]any{
		"title": t.title,
		"path":  prefix,
		"data":  all,
	}
}

// ParseData 解析路由
func (t *RouterData) ParseData(prefix string) []map[string]any {
	var all []map[string]any
	for _, datum := range t.data {
		all = append(all, map[string]any{
			"title":  datum.title,
			"name":   datum.name,
			"method": datum.method,
			"path":   prefix + datum.path,
		})
	}
	for _, item := range t.group {
		gpath := prefix + item.prefix
		data := item.ParseData(gpath)
		all = append(all, data...)
	}
	return all
}
