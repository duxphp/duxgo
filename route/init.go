package route

var Routes = map[string]*RouterData{}

// Add 添加路由
func Add(name string, route *RouterData) {
	Routes[name] = route
}

// Get 获取路由
func Get(name string) *RouterData {
	return Routes[name]
}
