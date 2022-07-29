package form

import (
	"github.com/duxphp/duxgo/core/ui/node"
)

// Number 文本输入框
type Number struct {
	before    map[string]any
	after     map[string]any
	step      any
	precision any
	min       any
	max       any
}

// NewNumber 创建文本
func NewNumber() *Number {
	return &Number{
		step:      1,
		precision: 0,
	}
}

// SetBeforeNumber 前置文本
func (a *Number) SetBeforeNumber(content any) *Number {
	a.before = map[string]any{
		"vSlot:prefix": "",
		"nodeName":     "span",
		"child":        content,
	}
	return a
}

// SetAfterNumber 后置文本
func (a *Number) SetAfterNumber(content any) *Number {
	a.after = map[string]any{
		"vSlot:suffix": "",
		"nodeName":     "span",
		"child":        content,
	}
	return a
}

// SetStep 补进设置
func (a *Number) SetStep(step any, precision int) *Number {
	a.step = step
	a.precision = precision
	return a
}

// SetLimit 限制设置
func (a *Number) SetLimit(min any, max any) *Number {
	a.min = min
	a.max = max
	return a
}

// GetValue 格式化值
func (a *Number) GetValue(value any, info map[string]any) any {
	return value
}

// SaveValue 保存数据
func (a *Number) SaveValue(value any, data map[string]any) any {
	return value
}

// Render 渲染
func (a *Number) Render(element node.IField) *node.TNode {

	var child []any

	if len(a.before) != 0 {
		child = append(child, a.before)
	}
	if len(a.after) != 0 {
		child = append(child, a.after)
	}

	ui := node.TNode{
		"nodeName":          "a-input-number",
		"vModel:modelValue": element.GetUIField(),
		"placeholder":       "请输入" + element.GetName(),
		"precision":         a.precision,
		"step":              a.step,
	}
	if len(child) != 0 {
		ui["child"] = child
	}
	if a.min != nil && a.max != nil {
		ui["min"] = a.min
		ui["max"] = a.max
	}
	return &ui
}
