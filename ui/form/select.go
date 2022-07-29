package form

import (
	"github.com/duxphp/duxgo/core/ui/node"
	"github.com/duxphp/duxgo/core/util/function"
	"github.com/samber/lo"
)

// Select 文本输入框
type Select struct {
	options     map[any]any
	placeholder string
	multi       bool
	url         string
	urlParams   map[string]any
	image       string
	desc        string
	icon        string
}

// NewSelect 创建文本
func NewSelect() *Select {
	return &Select{}

}

// SetMulti 设置多选
func (a *Select) SetMulti() *Select {
	a.multi = true
	return a
}

// SetOptions 设置选项
func (a *Select) SetOptions(options map[any]any) *Select {
	a.options = options
	return a
}

// SetPlaceholder 提示信息
func (a *Select) SetPlaceholder(content string) *Select {
	a.placeholder = content
	return a
}

// GetValue 格式化值
func (a *Select) GetValue(value any, info map[string]any) any {
	return value
}

// SaveValue 保存数据
func (a *Select) SaveValue(value any, data map[string]any) any {
	return value
}

// SetUrl 远程搜索
func (a *Select) SetUrl(url string, params map[string]any) *Select {
	a.url = url
	a.urlParams = params
	return a
}

func (a *Select) SetImage(field string, icons ...string) *Select {
	a.image = field
	if len(icons) > 0 {
		a.icon = icons[0]
	}
	return a
}
func (a *Select) SetDesc(field string) *Select {
	a.desc = field
	return a
}

// Render 渲染
func (a *Select) Render(element node.IField) *node.TNode {

	placeholder := a.placeholder

	if placeholder == "" {
		placeholder = "请输入" + element.GetName()
	}

	options := []map[string]any{}

	for value, label := range a.options {
		options = append(options, map[string]any{
			"label": label,
			"value": value,
		})
	}

	nParams := node.TNode{
		"placeholder": placeholder,
		"options":     options,
		"allowClear":  true,
		"multiple":    a.multi,
	}

	ui := node.TNode{
		"nodeName":     "app-select",
		"nParams":      nParams,
		"vModel:value": element.GetUIField(),
		"placeholder":  "请输入" + element.GetName(),
	}

	if a.url != "" {
		nParams["allowSearch"] = true
		nParams["filterOption"] = false
		ui["vBind:dataUrl"] = function.BuildUrl(a.url, a.urlParams, false)
	}

	if a.image != "" || a.desc != "" {
		mediaChild := []node.TNode{}
		if a.image != "" {
			mediaChild = append(mediaChild, node.TNode{
				"nodeName": "div",
				"class":    "flex-node",
				"child": node.TNode{
					"nodeName": "a-avatar",
					"size":     "34",
					"style": node.TNode{
						"backgroundColor": "#3370ff",
					},
					"child": []node.TNode{
						{
							"nodeName":  "img",
							"vIf":       "item.rowData.image != ''",
							"vBind:src": "item.rowData.image",
						},
						{
							"vIf":      "item.rowData.image == ''",
							"nodeName": lo.Ternary[string](a.icon == "", "IconUser", a.icon),
						},
					},
				},
			})
		}
		if a.desc != "" {
			mediaChild = append(mediaChild, node.TNode{
				"nodeName": "div",
				"class":    "flex-grow",
				"child": []node.TNode{
					{
						"nodeName": "div",
						"child":    "{{item.rowData.name}}",
					},
					{
						"nodeName": "div",
						"class":    "text-gray-400",
						"child":    "{{item.rowData.tel}}",
					},
				},
			})
		} else {
			mediaChild = append(mediaChild, node.TNode{
				"nodeName": "div",
				"class":    "flex-grow",
				"child":    "{{item.rowData.name}}",
			})
		}

		ui["vRender:optionRender:item"] = node.TNode{
			"nodeName": "div",
			"class":    "flex gap-2 py-2 items-center leading-none",
			"child":    mediaChild,
		}
	}

	return &ui
}
