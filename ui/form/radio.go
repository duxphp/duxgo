package form

import (
	"github.com/duxphp/duxgo/core/ui/node"
)

type RadioOptions struct {
	Key  any
	Name any
}

// Radio 文本输入框
type Radio struct {
	options []RadioOptions
}

// NewRadio 创建文本
func NewRadio() *Radio {
	return &Radio{}
}

// SetOptions 设置选项
func (a *Radio) SetOptions(options []RadioOptions) *Radio {
	a.options = options
	return a
}

// GetValue 格式化值
func (a *Radio) GetValue(value any, info map[string]any) any {
	if value == nil {
		value = a.options[0].Key
	}
	return value
}

// SaveValue 保存数据
func (a *Radio) SaveValue(value any, data map[string]any) any {
	return value
}

// Render 渲染
func (a *Radio) Render(element node.IField) *node.TNode {
	var options []map[string]any
	for _, item := range a.options {
		options = append(options, map[string]any{
			"nodeName": "a-radio",
			"child":    item.Name,
			"value":    item.Key,
		})
	}
	ui := node.TNode{
		"nodeName":          "a-radio-group",
		"child":             options,
		"vModel:modelValue": element.GetUIField(),
		"placeholder":       "请输入" + element.GetName(),
	}
	return &ui
}
