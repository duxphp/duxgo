package form

import (
	"github.com/duxphp/duxgo/ui/node"
)

// Text 文本输入框
type Text struct {
	num    int
	before map[string]any
	after  map[string]any
}

// NewText 创建文本
func NewText() *Text {
	return &Text{}
}

// SetNum 设置长度
func (a *Text) SetNum(num int) *Text {
	a.num = num
	return a
}

// SetBeforeText 前置文本
func (a *Text) SetBeforeText(content any) *Text {
	a.before = map[string]any{
		"vSlot:prepend": "",
		"nodeName":      "span",
		"child":         content,
	}
	return a
}

// SetAfterText 后置文本
func (a *Text) SetAfterText(content any) *Text {
	a.after = map[string]any{
		"vSlot:prepend": "",
		"nodeName":      "span",
		"child":         content,
	}
	return a
}

// GetValue 格式化值
func (a *Text) GetValue(value any, info map[string]any) any {
	return value
}

// SaveValue 保存数据
func (a *Text) SaveValue(value any, data map[string]any) any {
	return value
}

// Render 渲染
func (a *Text) Render(element node.IField) *node.TNode {

	var child []any

	if len(a.before) != 0 {
		child = append(child, a.before)
	}
	if len(a.after) != 0 {
		child = append(child, a.after)
	}

	ui := node.TNode{
		"nodeName":          "a-input",
		"vModel:modelValue": element.GetUIField(),
		"placeholder":       "请输入" + element.GetName(),
	}
	if len(child) != 0 {
		ui["child"] = child
	}

	return &ui
}
