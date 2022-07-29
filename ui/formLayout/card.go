package formLayout

import (
	form2 "github.com/duxphp/duxgo/ui/form"
	"github.com/duxphp/duxgo/ui/node"
)

type Card struct {
	form   *form2.Form
	data   map[string]any
	dialog bool
}

// NewCard 卡片布局
func NewCard() *Card {
	return &Card{}
}

// SetData 设置数据
func (a *Card) SetData(data map[string]any) {
	a.data = data
}

// SetDialog 设置弹窗
func (a *Card) SetDialog(dialog bool) {
	a.dialog = dialog
}

// Column 列元素
func (a *Card) Column(callback func(form *form2.Form), opt ...any) {
	formUI := form2.NewForm()
	formUI.SetData(a.data)
	formUI.SetDialog(a.dialog)
	a.form = formUI
	callback(a.form)
}

// Form 获取表单
func (a *Card) Form(index ...int) *form2.Form {
	return a.form
}

// Expand 展开元素
func (a *Card) Expand() []*form2.Element {
	return a.form.ExpandElement()
}

// Render 渲染
func (a *Card) Render() *node.TNode {
	element := a.form.RenderElement()
	ui := node.TNode{
		"nodeName": "div",
		"class":    "mb-4 bg-white dark:bg-blackgray-4 rounded shadow p-7 pb-2",
		"child":    element,
	}
	return &ui
}
