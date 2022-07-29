package form

import (
	"github.com/duxphp/duxgo/core/ui/node"
)

// Editor 编辑器
type Editor struct {
	num    int
	before map[string]any
	after  map[string]any
}

// NewEditor 创建文本
func NewEditor() *Editor {
	return &Editor{}
}

// GetValue 格式化值
func (a *Editor) GetValue(value any, info map[string]any) any {
	return value
}

// SaveValue 保存数据
func (a *Editor) SaveValue(value any, data map[string]any) any {
	return value
}

// Render 渲染
func (a *Editor) Render(element node.IField) *node.TNode {
	ui := node.TNode{
		"nodeName":     "app-editor",
		"vModel:value": element.GetUIField(),
	}
	return &ui
}
