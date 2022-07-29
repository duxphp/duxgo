package formLayout

import (
	"github.com/duxphp/duxgo/core/ui/form"
	"github.com/duxphp/duxgo/core/ui/node"
)

type TabArgs struct {
	Name  string
	Order uint
	Title string
	Desc  string
}

type TabColumn struct {
	form   *form.Form
	params TabArgs
}

type Tab struct {
	data    map[string]any
	dialog  bool
	columns []TabColumn
}

// NewTab 切换布局
func NewTab() *Tab {
	form.NewForm()
	return &Tab{}
}

// SetData 设置数据
func (a *Tab) SetData(data map[string]any) {
	a.data = data
}

// SetDialog 设置弹窗
func (a *Tab) SetDialog(dialog bool) {
	a.dialog = dialog
}

// Column 列元素
func (a *Tab) Column(callback func(form *form.Form), opt ...any) {
	formUI := form.NewForm()
	formUI.SetData(a.data)
	formUI.SetDialog(a.dialog)
	params := TabArgs{}
	if len(opt) > 0 {
		params = opt[0].(TabArgs)
	}
	a.columns = append(a.columns, TabColumn{
		form:   formUI,
		params: params,
	})
	callback(formUI)
}

// Form 获取表单
func (a *Tab) Form(index ...int) *form.Form {
	i := 0
	if len(index) > 0 {
		i = index[0]
	}
	return a.columns[i].form
}

// Expand 展开元素
func (a *Tab) Expand() []*form.Element {
	var element []*form.Element
	for _, column := range a.columns {
		element = append(element, column.form.ExpandElement()...)
	}
	return element
}

// Render 渲染
func (a *Tab) Render() *node.TNode {
	var children []node.TNode
	for key, column := range a.columns {
		element := column.form.RenderElement()
		var child []*node.TNode
		if column.params.Title != "" {
			child = append(child, &node.TNode{
				"nodeName": "div",
				"class":    "py-4 flex flex-col gap-2",
				"child": []node.TNode{
					{
						"nodeName": "div",
						"class":    "text-xl",
						"child":    column.params.Title,
					},
					{
						"nodeName": "div",
						"class":    "text-gray-500",
						"child":    column.params.Desc,
					},
				},
			})
		}
		child = append(child, element...)
		children = append(children, node.TNode{
			"nodeName": "a-tab-pane",
			"title":    column.params.Name,
			"key":      key,
			"class":    "border-t border-gray-200 dark:border-blackgray-1 px-3 pt-4 pb-0",
			"child":    child,
		})

	}

	ui := node.TNode{
		"nodeName": "a-tabs",
		"class":    "mb-4 bg-white dark:bg-blackgray-4 rounded shadow p-4 pb-1",
		"type":     "rounded",
		"child":    children,
	}
	return &ui
}
