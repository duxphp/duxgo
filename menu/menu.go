package menu

import (
	"sort"
)

// MenuData 菜单应用结构
type MenuData struct {
	App      string `json:"app"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Icon     string `json:"icon"`
	Title    string `json:"title"`
	Hidden   bool   `json:"hidden"`
	Order    int    `json:"order"`
	Data     []*MenuData
	PushData map[string]*MenuData
}

// New 新建菜单
func New() *MenuData {
	return &MenuData{
		PushData: map[string]*MenuData{},
	}
}

// Add 添加菜单
func (t *MenuData) Add(data *MenuData) *MenuData {
	t.Data = append(t.Data, data)
	return data
}

// Push 追加菜单
func (t *MenuData) Push(app string) *MenuData {
	data := &MenuData{App: app}
	t.PushData[data.App] = data
	return data
}

// Group 添加菜单组
func (t *MenuData) Group(name string) *MenuData {
	data := &MenuData{
		Name: name,
	}
	t.Data = append(t.Data, data)
	return data
}

// Item 添加条目
func (t *MenuData) Item(name string, url string, order int) {
	data := &MenuData{
		Name:  name,
		Url:   url,
		Order: order,
	}
	t.Data = append(t.Data, data)
}

// Get 获取菜单
func (t *MenuData) Get() map[string]any {
	// 重置菜单
	var menu []map[string]any
	for _, appData := range t.Data {
		// 合并追加菜单
		if t.PushData[appData.App] != nil {

			for _, datum := range t.PushData[appData.App].Data {
				appData.Data = append(appData.Data, datum)
			}
		}
		// 重置分组菜单
		var group []map[string]any
		for _, groupData := range appData.Data {
			// 重置子菜单
			var list []map[string]any
			for _, items := range groupData.Data {
				list = append(list, map[string]any{
					"name":  items.Name,
					"url":   items.Url,
					"order": items.Order,
				})
			}
			sort.Slice(list, func(i, j int) bool {
				return list[i]["order"].(int) < list[j]["order"].(int)
			})
			group = append(group, map[string]any{
				"name":  groupData.Name,
				"order": groupData.Order,
				"title": groupData.Title,
				"menu":  list,
			})
		}
		sort.Slice(group, func(i, j int) bool {
			return group[i]["order"].(int) < group[j]["order"].(int)
		})
		menu = append(menu, map[string]any{
			"name":  appData.Name,
			"icon":  appData.Icon,
			"order": appData.Order,
			"url":   appData.Url,
			"menu":  group,
		})
	}
	sort.Slice(menu, func(i, j int) bool {
		return menu[i]["order"].(int) < menu[j]["order"].(int)
	})

	return map[string]any{
		"list": menu,
		"app":  []any{},
	}
}
