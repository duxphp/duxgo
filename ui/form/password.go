package form

import (
	"github.com/duxphp/duxgo/core/ui/node"
	"github.com/duxphp/duxgo/core/util/function"
	"github.com/spf13/cast"
)

// Password 文本输入框
type Password struct {
}

// NewPassword 创建文本
func NewPassword() *Password {
	return &Password{}
}

// GetValue 格式化值
func (a *Password) GetValue(value any, info map[string]any) any {
	return value
}

// SaveValue 保存数据
func (a *Password) SaveValue(value any, data map[string]any) any {
	val := cast.ToString(value)
	if val == "" {
		return ""
	} else {
		return function.HashEncode([]byte(val))
	}
}

// Render 渲染
func (a *Password) Render(element node.IField) *node.TNode {
	ui := node.TNode{
		"nodeName":          "a-input-password",
		"vModel:modelValue": element.GetUIField(),
		"placeholder":       "请输入" + element.GetName(),
		"allowClear":        true,
	}

	return &ui
}
