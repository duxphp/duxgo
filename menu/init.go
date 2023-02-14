package menu

var Menus = map[string]*MenuData{}

// Add 添加菜单
func Add(name string, route *MenuData) {
	Menus[name] = route
}

// Get 获取菜单
func Get(name string) *MenuData {
	return Menus[name]
}
