package widget

import (
	"github.com/duxphp/duxgo/core/ui/node"
	"github.com/duxphp/duxgo/core/util/function"
	"github.com/samber/lo"
)

type Menu struct {
	Url    string
	Event  string
	Name   string
	Type   string
	Key    string
	Method string
}

type treeList struct {
	key       string            // 当前选中Key
	field     string            // 查询字段
	data      []map[string]any  // 树形数据
	url       string            // 数据 Url
	sortUrl   string            // 排序 Url
	labelNode node.TNode        // 标签节点
	search    bool              // 搜索状态
	keywords  []string          // 搜索字段
	filter    string            // js过滤条件
	fieldMaps map[string]string // 字段映射
	event     string            // 事件名称
	menus     []Menu            // 菜单数据
}

func NewTreeList(key string, field string) *treeList {
	return &treeList{
		key:   key,
		field: field,
	}
}

// SetData 设置数据
func (a *treeList) SetData(data []map[string]any) *treeList {
	a.data = data
	return a
}

// SetUrl 设置数据Url
func (a *treeList) SetUrl(url string) *treeList {
	a.url = url
	return a
}

// SetSortUrl 设置排序Url
func (a *treeList) SetSortUrl(url string) *treeList {
	a.sortUrl = url
	return a
}

// SetLabel 设置标签节点
func (a *treeList) SetLabel(node node.TNode) *treeList {
	a.labelNode = node
	return a
}

// SetSearch 设置搜索
func (a *treeList) SetSearch(search bool, keywords []string) *treeList {
	a.search = search
	a.keywords = keywords
	return a
}

// SetFilter 设置js过滤
func (a *treeList) SetFilter(filter string) *treeList {
	a.filter = filter
	return a
}

// SetFieldMaps 字段映射
func (a *treeList) SetFieldMaps(field map[string]string) *treeList {
	a.fieldMaps = field
	return a
}

// SetEvent 设置事件
func (a *treeList) SetEvent(event string) *treeList {
	a.event = function.Md5(event)
	return a
}

// SetMenu 设置菜单
func (a *treeList) SetMenu(menu []Menu) *treeList {
	a.menus = menu
	return a
}

// Render 渲染UI
func (a *treeList) Render() node.TNode {

	if a.url != "" && a.event == "" {
		a.event = function.Md5(a.url)
	}
	element := node.TNode{
		"nodeName":         "widget-tree",
		"treeData":         a.data,
		"url":              a.url,
		"sortUrl":          a.sortUrl,
		"search":           a.search,
		"keywords":         a.keywords,
		"requestEventName": a.event,
		"vBind:filter":     a.filter,
		"iconColor":        []string{"blue", "cyan", "green", "orange", "red", "purple"},
		"vModel:value":     "data.filter['" + a.field + "']",
	}

	if a.fieldMaps != nil {
		element["fieldNames"] = a.fieldMaps
	}

	if a.labelNode != nil {
		element["child"] = node.TNode{
			"nodeName":    "span",
			"vSlot:label": "item",
			"child":       a.labelNode,
		}
	}

	var menus []map[string]any
	for _, menu := range a.menus {
		url := menu.Url
		event := menu.Event
		tmp := map[string]any{
			"text": menu.Name,
			"key":  menu.Key,
		}

		if event != "" {
			tmp["event"] = event
		} else {
			switch menu.Type {
			case "dialog":
				tmp["event"] = lo.Ternary[string](url != "", `window.router.dialog(`+url+`)`, `window.dialog.alert({content: '未定义链接数据'})`)
			case "ajax":
				tmp["event"] = lo.Ternary[string](url != "", `window.router.ajax(`+url+`, {_method: '`+menu.Method+`', _title: '确认进行`+menu.Name+`操作？'})`, `window.dialog.alert({content: '未定义链接数据'})`)
			default:
				tmp["event"] = lo.Ternary[string](url != "", `window.router.push(`+url+`)`, `window.dialog.alert({content: '未定义链接数据'})`)
			}

		}
		menus = append(menus, tmp)
	}
	if len(menus) > 0 {
		element["contextMenus"] = menus
	}
	return element
}
