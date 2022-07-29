package column

import (
	"github.com/duxphp/duxgo/ui/node"
	"github.com/duxphp/duxgo/ui/table"
	"github.com/duxphp/duxgo/ui/widget"
)

// Link 链接
type Link struct {
	links []*widget.Link
}

// NewLink 创建文本
func NewLink() *Link {
	return &Link{}
}

// AddUrl 添加描述
func (a *Link) AddUrl(name string, url string, params ...map[string]any) *widget.Link {
	link := widget.NewLink(name, url, params...)
	a.links = append(a.links, link)
	return link
}

// Render 渲染
func (a *Link) Render(element *table.Column) node.TNode {

	var el []node.TNode
	for _, item := range a.links {
		el = append(el, node.TNode{
			"nodeName": "span",
			"child":    item.SetModel("rowData.record", element.GetFields()).Render(),
		})
	}
	return node.TNode{
		"nodeName": "div",
		"class":    "inline-flex gap-2",
		"child":    el,
	}
}
