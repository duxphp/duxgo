package form

import (
	"github.com/duxphp/duxgo/ui/node"
	"github.com/duxphp/duxgo/util/function"
	"github.com/spf13/cast"
)

// IElement 字段UI接口
type IElement interface {
	Render(element node.IField) *node.TNode
	GetValue(value any, info map[string]any) any
	SaveValue(value any, data map[string]any) any
}

// ILayout 布局UI接口
type ILayout interface {
	SetData(map[string]any)
	SetDialog(bool)
	Column(callback func(form *Form), opt ...any)
	Form(index ...int) *Form
	Render() *node.TNode
	Expand() []*Element
}

// Element 字段结构
type Element struct {
	Name     string
	Field    string
	HasAs    string
	HasModel any
	HasKey   string
	Default  any
	Help     any
	HelpLine bool
	Must     bool
	Verify   []map[string]string
	Format   func(value any) any
	Value    *any
	UI       IElement
	Layout   ILayout
}

// SetDefault 设置默认值
func (a *Element) SetDefault(value any) *Element {
	a.Default = value
	return a
}

// SetUI 设置UI
func (a *Element) SetUI(ui IElement) *Element {
	a.UI = ui
	return a
}

// SetHas 设置关联模型
func (a *Element) SetHas(as string, model any, key ...string) *Element {
	a.HasAs = as
	a.HasModel = model
	if len(key) > 0 {
		a.HasKey = key[0]
	} else {
		a.HasKey = "id"
	}
	return a
}

// SetMust 设置必填
func (a *Element) SetMust(status bool) *Element {
	a.Must = status
	a.SetVerify("required")
	return a
}

// SetVerify 设置必填
func (a *Element) SetVerify(role string, message ...string) *Element {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	a.Verify = append(a.Verify, map[string]string{
		"role":    role,
		"message": msg,
	})
	return a
}

// SetValue 设置覆盖值
func (a *Element) SetValue(val any) *Element {
	a.Value = &val
	return a
}

// SetHelp 设置帮助信息
func (a *Element) SetHelp(help any, line ...bool) *Element {
	a.Help = help
	if len(line) > 0 {
		a.HelpLine = line[0]
	} else {
		a.HelpLine = true
	}
	return a
}

// SaveFormat 保存格式化
func (a *Element) SaveFormat(callback func(value any) any) *Element {
	a.Format = callback
	return a
}

// SetType 设置类型
func (a *Element) SetType() *Element {
	return a
}

// GetName 获取名称
func (a *Element) GetName() string {
	return a.Name
}

// GetData 获取值
func (a *Element) GetData(info map[string]any) any {
	var content any
	if a.HasAs != "" {
		// 多对多数据
		hasMap := info[function.LcFirst(a.HasAs)]
		var ids []any
		for _, item := range cast.ToSlice(hasMap) {
			tmp := cast.ToStringMap(item)
			ids = append(ids, tmp[a.HasKey])
		}
		content = ids
	} else {
		// 普通数据
		content = info[a.Field]
	}
	if a.Value != nil {
		content = a.Value
	}
	if a.UI != nil {
		content = a.UI.GetValue(content, info)
	}
	if content == nil && a.Default != nil {
		content = a.Default
	}
	return content
}

// GetUIField 获取字段名
func (a *Element) GetUIField(field ...string) string {
	content := a.Field
	if len(field) > 0 {
		content = field[0]
	}
	return "data['" + content + "']"
}

// SaveData 保存数据处理
func (a *Element) SaveData(value any, data map[string]any) any {
	if a.UI != nil {
		value = a.UI.SaveValue(value, data)
	}
	return value
}

// Render 渲染UI
func (a *Element) Render() *node.TNode {

	ui := a.UI.Render(a)

	var helpNode node.TNode

	if a.HelpLine {
		helpNode = node.TNode{
			"nodeName":   "div",
			"vSlot:help": "",
			"class":      "mb-2",
			"child":      a.Help,
		}
	} else {
		helpNode = node.TNode{
			"nodeName": "div",
			"class":    "ml-2",
			"child":    a.Help,
		}
	}

	nodeEl := node.TNode{
		"nodeName": "a-form-item",
		"label":    a.Name,
		"field":    a.Field,
		"child": []node.TNode{
			*ui,
			helpNode,
		},
	}

	if a.Must {
		nodeEl["rules"] = node.TNode{
			"required": true,
			"message":  "请填写" + a.Name,
		}
	}
	return &nodeEl
}
