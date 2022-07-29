package formLayout

import (
	form2 "github.com/duxphp/duxgo/ui/form"
	"github.com/duxphp/duxgo/ui/node"
	"github.com/spf13/cast"
)

type RowColumn struct {
	form *form2.Form
	span int
}

type Row struct {
	data    map[string]any
	dialog  bool
	columns []RowColumn
}

// NewRow 切换布局
func NewRow() *Row {
	form2.NewForm()
	return &Row{}
}

// SetData 设置数据
func (a *Row) SetData(data map[string]any) {
	a.data = data
}

// SetDialog 设置弹窗
func (a *Row) SetDialog(dialog bool) {
	a.dialog = dialog
}

// Column 列元素
func (a *Row) Column(callback func(form *form2.Form), opt ...any) {
	formUI := form2.NewForm()
	formUI.SetData(a.data)
	formUI.SetDialog(a.dialog)
	a.columns = append(a.columns, RowColumn{
		form: formUI,
		span: cast.ToInt(opt[0]),
	})
	callback(formUI)
}

// Form 获取表单
func (a *Row) Form(index ...int) *form2.Form {
	i := 0
	if len(index) > 0 {
		i = index[0]
	}
	return a.columns[i].form
}

// Expand 展开元素
func (a *Row) Expand() []*form2.Element {
	var element []*form2.Element
	for _, column := range a.columns {
		element = append(element, column.form.ExpandElement()...)
	}
	return element
}

// Render 渲染
func (a *Row) Render() *node.TNode {
	var children []node.TNode
	for _, column := range a.columns {
		element := column.form.RenderElement()
		children = append(children, node.TNode{
			"nodeName": "a-col",
			"span":     column.span,
			"child":    element,
		})

	}

	ui := node.TNode{
		"nodeName": "a-row",
		"gutter":   "24",
		"child":    children,
	}
	return &ui
}
