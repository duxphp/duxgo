package permission

var Permissions = map[string]*PermissionData{}

// Add 添加权限
func Add(name string, route *PermissionData) {
	Permissions[name] = route
}

// Get 获取权限
func Get(name string) *PermissionData {
	return Permissions[name]
}
