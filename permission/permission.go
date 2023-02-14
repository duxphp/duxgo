package menu

import (
	"github.com/samber/lo"
	"sort"
)

// PermissionData 权限应用结构
type PermissionData struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Order int    `json:"order"`
	Data  []*PermissionData
}

// New 新建菜单
func New() *PermissionData {
	return &PermissionData{}
}

// Group 添加权限组
func (t *PermissionData) Group(name string, label string, order int) *PermissionData {
	data := &PermissionData{
		Name:  name,
		Label: label,
		Order: order,
	}
	t.Data = append(t.Data, data)
	return data
}

// Add 添加权限
func (t *PermissionData) Add(label string, name string) {
	data := &PermissionData{
		Name:  name,
		Label: t.Label + "." + label,
	}
	t.Data = append(t.Data, data)
}

// Get 获取权限
func (t *PermissionData) Get() []map[string]any {
	data := lo.Map[*PermissionData, map[string]any](t.Data, func(group *PermissionData, index int) map[string]any {
		list := lo.Map[*PermissionData, map[string]any](group.Data, func(item *PermissionData, index int) map[string]any {
			return map[string]any{
				"name":  item.Name,
				"label": item.Label,
			}
		})
		return map[string]any{
			"label":    group.Label,
			"order":    group.Order,
			"name":     group.Name,
			"children": list,
		}
	})
	sort.Slice(data, func(i, j int) bool {
		return data[i]["order"].(int) < data[j]["order"].(int)
	})
	return data
}
