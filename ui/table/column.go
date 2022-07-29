package table

import (
	"github.com/duxphp/duxgo/core/ui/node"
	"strings"
)

// IColumn 字段UI接口
type IColumn interface {
	Render(element *Column) node.TNode
}

// Column 字段结构
type Column struct {
	Name   string
	Field  string
	Node   node.TNode
	Width  uint
	Format func(value any, data map[string]any) any
	Fields []string
	UI     IColumn
	Sort   bool
}

// SetUI 设置UI
func (a *Column) SetUI(ui IColumn) *Column {
	a.UI = ui
	return a
}

// SetWidth 设置宽度
func (a *Column) SetWidth(width uint) *Column {
	a.Width = width
	return a
}

// SetSort 设置排序
func (a *Column) SetSort() *Column {
	a.Sort = true
	return a
}

// SetNode 自定义元素
func (a *Column) SetNode(node node.TNode) *Column {
	a.Node = node
	return a
}

// DataFormat 数据格式化
func (a *Column) DataFormat(callback func(value any, data map[string]any) any) *Column {
	a.Format = callback
	return a
}

// GetName 获取名称
func (a *Column) GetName() string {
	return a.Name
}

// GetUIField 获取模型字段名
func (a *Column) GetUIField(field ...string) string {
	content := a.Field
	if len(field) > 0 {
		content = field[0]
	}
	newField := strings.Replace(content, ".", "_", -1)
	return "rowData.record['" + newField + "']"
}

// GetFields 获取所有字段
func (a *Column) GetFields() []string {
	return a.Fields
}

// Render 渲染UI
func (a *Column) Render(fields []string) node.TNode {
	a.Fields = fields
	el := node.TNode{
		"title":     a.Name,
		"dataIndex": strings.Replace(a.Field, ".", "_", -1),
		"width":     a.Width,
	}
	if a.Sort {
		el["vBind:sortable"] = "colSortable"
	}
	render := a.Node
	if render == nil && a.UI != nil {
		render = a.UI.Render(a)
	}
	if render != nil {
		el["render:rowData, rowIndex"] = render
	}
	return el
}
