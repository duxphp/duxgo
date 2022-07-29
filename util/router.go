package util

import (
	"github.com/labstack/echo/v4"
)

// RouterData 路由结构
type RouterData struct {
	name       string
	prefix     string
	permission bool
	data       []*RouterItem
	group      []*RouterData
	router     *echo.Group
}

// RouterItem 路由结构
type RouterItem struct {
	name   string
	method string
	path   string
	as     string
}

// NewRouter 新建资源路由
func NewRouter(router *echo.Group) *RouterData {
	return &RouterData{
		router: router,
	}
}

// Group 路由分组
func (t *RouterData) Group(prefix string, name string, middle ...echo.MiddlewareFunc) *RouterData {
	group := &RouterData{
		name:   name,
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
func (t *RouterData) Get(path string, handler echo.HandlerFunc, name string, as ...string) *echo.Route {
	return t.Add(echo.GET, path, handler, name, as...)
}

// Head 路由
func (t *RouterData) Head(path string, handler echo.HandlerFunc, name string, as ...string) *echo.Route {
	return t.Add(echo.HEAD, path, handler, name, as...)
}

// Post 路由
func (t *RouterData) Post(path string, handler echo.HandlerFunc, name string, as ...string) *echo.Route {
	return t.Add(echo.POST, path, handler, name, as...)
}

// Put 路由
func (t *RouterData) Put(path string, handler echo.HandlerFunc, name string, as ...string) *echo.Route {
	return t.Add(echo.PUT, path, handler, name, as...)
}

// Delete 路由
func (t *RouterData) Delete(path string, handler echo.HandlerFunc, name string, as ...string) *echo.Route {
	return t.Add(echo.DELETE, path, handler, name, as...)
}

// Connect 路由
func (t *RouterData) Connect(path string, handler echo.HandlerFunc, name string, as ...string) *echo.Route {
	return t.Add(echo.CONNECT, path, handler, name, as...)
}

// Options 路由
func (t *RouterData) Options(path string, handler echo.HandlerFunc, name string, as ...string) *echo.Route {
	return t.Add(echo.OPTIONS, path, handler, name, as...)
}

// Trace 路由
func (t *RouterData) Trace(path string, handler echo.HandlerFunc, name string, as ...string) *echo.Route {
	return t.Add(echo.TRACE, path, handler, name, as...)
}

// Patch 路由
func (t *RouterData) Patch(path string, handler echo.HandlerFunc, name string, as ...string) *echo.Route {
	return t.Add(echo.PATCH, path, handler, name, as...)
}

// Add 添加路由资源
func (t *RouterData) Add(method string, path string, handler echo.HandlerFunc, name string, as ...string) *echo.Route {
	item := RouterItem{
		name:   name,
		method: method,
		path:   path,
		as:     "",
	}
	if len(as) > 0 {
		item.as = as[0]
	}
	t.data = append(t.data, &item)
	r := t.router.Add(method, path, handler)
	r.Name = item.as
	return r
}

// ParseTree 解析路由为树形
func (t *RouterData) ParseTree(prefix string) any {
	var all []any
	for _, datum := range t.data {
		all = append(all, map[string]any{
			"name":       datum.name,
			"method":     datum.method,
			"path":       prefix + datum.path,
			"permission": t.permission,
		})
	}
	for _, item := range t.group {
		gpath := prefix + item.prefix
		all = append(all, item.ParseTree(gpath))
	}
	if t.name == "" {
		return all
	}
	return map[string]any{
		"name":       t.name,
		"permission": t.permission,
		"path":       prefix,
		"data":       all,
	}
}

// ParseData 解析路由
func (t *RouterData) ParseData(prefix string) []map[string]any {
	var all []map[string]any
	for _, datum := range t.data {
		all = append(all, map[string]any{
			"name":       datum.name,
			"method":     datum.method,
			"path":       prefix + datum.path,
			"permission": t.permission,
		})
	}
	for _, item := range t.group {
		gpath := prefix + item.prefix
		data := item.ParseData(gpath)
		all = append(all, data...)
	}
	return all
}
