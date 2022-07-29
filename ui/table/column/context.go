package column

import (
	"github.com/duxphp/duxgo/core/ui/node"
	"github.com/duxphp/duxgo/core/ui/table"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"strings"
)

// Context 图文展示
type Context struct {
	desc  []map[string]any
	image []map[string]any
}

// NewContext 创建文本
func NewContext() *Context {
	return &Context{
		desc:  make([]map[string]any, 0),
		image: make([]map[string]any, 0),
	}
}

// AddDesc 添加描述
func (a *Context) AddDesc(field string) *Context {
	newField := strings.Replace(field, ".", "_", -1)
	a.desc = append(a.desc, map[string]any{
		"name": newField,
	})
	return a
}

// AddImage 添加图片
func (a *Context) AddImage(field string, size int, icons ...string) *Context {
	newField := strings.Replace(field, ".", "_", -1)
	icon := ""
	if len(icons) > 0 {
		icon = icons[0]
	}
	a.image = append(a.image, map[string]any{
		"name": newField,
		"size": size,
		"icon": icon,
	})
	return a
}

// Render 渲染
func (a *Context) Render(element *table.Column) node.TNode {

	var el []node.TNode

	for _, item := range a.image {
		field := element.GetUIField(item["name"].(string))
		el = append(el, node.TNode{
			"nodeName": "div",
			"class":    "flex-none",
			"child": node.TNode{
				"nodeName": "a-avatar",
				"size":     item["size"],
				"style": node.TNode{
					"backgroundColor": "#3370ff",
				},
				"child": []node.TNode{
					{
						"nodeName":  "img",
						"vIf":       field + " != ''",
						"vBind:src": field,
					},
					{
						"vIf":      field + " == ''",
						"nodeName": lo.Ternary[string](item["icon"] == "", "iconImage", cast.ToString(item["icon"])),
					},
				},
			},
		})
	}

	var content []node.TNode
	content = append(content, node.TNode{
		"nodeName": "div",
		"child":    "{{" + element.GetUIField() + " || '-'}}",
	})
	for _, item := range a.desc {
		content = append(content, node.TNode{
			"nodeName": "div",
			"class":    "text-gray-500",
			"child":    "{{" + element.GetUIField(item["name"].(string)) + " || '-'}}",
		})
	}

	el = append(el, node.TNode{
		"nodeName": "div",
		"class":    "flex-grow",
		"child":    content,
	})

	return node.TNode{
		"nodeName": "div",
		"class":    "flex flex-row gap-2 items-center",
		"child":    el,
	}
}
