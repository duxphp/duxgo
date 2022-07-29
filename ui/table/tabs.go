package table

import (
	"github.com/duxphp/duxgo/ui/node"
	"github.com/jianfengye/collection"
	"gorm.io/gorm"
)

// Tab 筛选结构
type Tab struct {
	Name    string
	Icon    string
	Where   func(*gorm.DB)
	Collect func(*collection.ICollection)
}

// SetIcon 设置图标
func (a *Tab) SetIcon(icon string) *Tab {
	a.Icon = icon
	return a
}

// SetWhere 设置条件
func (a *Tab) SetWhere(where func(*gorm.DB)) *Tab {
	a.Where = where
	return a
}

// SetCollect 设置集合
func (a *Tab) SetCollect(where func(*collection.ICollection)) *Tab {
	a.Collect = where
	return a
}

// Render 渲染UI
func (a *Tab) Render(index int) node.TNode {
	icon := node.TNode{}
	if a.Icon != "" {
		icon = node.TNode{
			"nodeName": a.Icon,
		}
	}
	return node.TNode{
		"nodeName": "a-radio",
		"value":    index,
		"child": []node.TNode{
			icon,
			{
				"nodeName": "span",
				"child":    " " + a.Name,
			},
		},
	}
}
