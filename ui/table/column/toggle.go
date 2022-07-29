package column

import (
	"github.com/duxphp/duxgo/core/ui/node"
	"github.com/duxphp/duxgo/core/ui/table"
	"github.com/duxphp/duxgo/core/util/function"
)

// Toggle 快速切换
type Toggle struct {
	url    string
	field  string
	params map[string]any
}

// NewToggle 创建组件
func NewToggle(field string) *Toggle {
	return &Toggle{
		field: field,
	}
}

// SetUrl 设置Url
func (a *Toggle) SetUrl(url string, params map[string]any) *Toggle {
	a.url = url
	a.params = params
	return a
}

// Render 渲染
func (a *Toggle) Render(element *table.Column) node.TNode {
	modelField := element.Field
	url := function.BuildUrl(a.url, a.params, false, "rowData.record", element.Fields)
	return node.TNode{
		"nodeName":          "a-switch",
		"vModel:modelValue": element.GetUIField(),
		"vOn:change":        `rowData.record['` + modelField + `'] = $event, editValue(` + url + `, {'field': '` + a.field + `', '` + a.field + `': rowData.record['` + modelField + `']})`,
	}
}
