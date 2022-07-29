package column

import (
	"github.com/duxphp/duxgo/ui/node"
	"github.com/duxphp/duxgo/ui/table"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// Image 图片
type Image struct {
	size uint
	icon string
}

// NewImage 创建组件
func NewImage(size uint, icons ...string) *Image {
	icon := ""
	if len(icons) > 0 {
		icon = icons[0]
	}
	return &Image{
		size: size,
		icon: icon,
	}
}

// Render 渲染
func (a *Image) Render(element *table.Column) node.TNode {
	modelField := element.Field
	return node.TNode{
		"nodeName": "a-avatar",
		"size":     a.size,
		"style": node.TNode{
			"backgroundColor": "#3370ff",
		},
		"child": []node.TNode{
			{
				"nodeName":  "img",
				"vIf":       `rowData.record['` + modelField + `']` + " != ''",
				"vBind:src": `rowData.record['` + modelField + `']`,
			},
			{
				"vIf":      `rowData.record['` + modelField + `']` + " == ''",
				"nodeName": lo.Ternary[string](a.icon == "", "iconImage", cast.ToString(a.icon)),
			},
		},
	}
}
