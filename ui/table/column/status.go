package column

import (
	"github.com/duxphp/duxgo/core/ui/node"
	"github.com/duxphp/duxgo/core/ui/table"
	"fmt"
)

// Status 状态展示
type Status struct {
	maps   map[any]string
	colors map[any]string
	types  string
}

// NewStatus 创建状态
func NewStatus(maps *map[any]string, colors *map[any]string) *Status {
	return &Status{
		maps:   *maps,
		colors: *colors,
	}
}

// SetText 设置文字状态
func (a *Status) SetText() *Status {
	a.types = "text"
	return a
}

// Render 渲染
func (a *Status) Render(element *table.Column) node.TNode {

	data := map[any]map[string]any{}

	for k, v := range a.maps {
		data[k] = map[string]any{
			"name":  v,
			"color": a.colors[k],
		}
	}

	var el []node.TNode
	field := element.GetUIField()
	for key, item := range data {

		isNum := false
		switch key.(type) {
		case int:
			isNum = true
		}
		vIf := fmt.Sprintf(`%v == '%v'`, field, key)
		if isNum {
			vIf = fmt.Sprintf(`%v == %v`, field, key)
		}
		el = append(el, node.TNode{
			"nodeName": "a-tag",
			"color":    item["color"],
			"size":     "medium",
			"child":    item["name"],
			"vIf":      vIf,
		})
	}
	return node.TNode{
		"nodeName": "div",
		"class":    "flex flex-row gap-2",
		"child":    el,
	}
}
