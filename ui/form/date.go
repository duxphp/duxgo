package form

import (
	"github.com/duxphp/duxgo/core/ui/node"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"gorm.io/datatypes"
)

type DateType string

const (
	DateYear  DateType = "year"
	DateMonth DateType = "month"
)

// Date 日期输入框
type Date struct {
	mode DateType
}

// NewDate 创建日期
func NewDate() *Date {
	return &Date{}
}

// Type 日期模式
func (a *Date) Type(mode DateType) *Date {
	a.mode = mode
	return a
}

// GetValue 格式化值
func (a *Date) GetValue(value any, info map[string]any) any {
	val := cast.ToTime(value)
	switch a.mode {
	case "year":
		return lo.Ternary[string](val.IsZero(), "", val.Format("2006"))
	case "month":
		return lo.Ternary[string](val.IsZero(), "", val.Format("2006-01"))
	default:
		return lo.Ternary[string](val.IsZero(), "", val.Format("2006-01-02"))
	}
}

// SaveValue 保存数据
func (a *Date) SaveValue(value any, data map[string]any) any {
	val := cast.ToString(value)
	if val == "" {
		return nil
	} else {
		return datatypes.Date(cast.ToTime(val))
	}
}

// Render 渲染
func (a *Date) Render(element node.IField) *node.TNode {
	ui := node.TNode{
		"nodeName":          "a-date-picker",
		"allowClear":        true,
		"vModel:modelValue": element.GetUIField(),
		"placeholder":       "请输入" + element.GetName(),
	}
	if a.mode != "" {
		ui["nodeName"] = "a-" + a.mode + "-picker"
	}
	return &ui
}
