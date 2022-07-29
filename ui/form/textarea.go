package form

import (
	"github.com/duxphp/duxgo/core/ui/node"
)

// Textarea 多行文本框
type Textarea struct {
	num    int
	before map[string]any
	after  map[string]any
}

// NewTextarea 创建多行文本
func NewTextarea() *Textarea {
	return &Textarea{}
}

// GetValue 格式化值
func (a *Textarea) GetValue(value any, info map[string]any) any {
	return value
}

// SaveValue 保存数据
func (a *Textarea) SaveValue(value any, data map[string]any) any {
	return value
}

// Render 渲染
func (a *Textarea) Render(element node.IField) *node.TNode {

	ui := node.TNode{
		"nodeName":          "a-textarea",
		"vModel:modelValue": element.GetUIField(),
		"placeholder":       "请输入" + element.GetName(),
		"allowClear":        true,
		"showWordLimit":     true,
	}

	return &ui
}
