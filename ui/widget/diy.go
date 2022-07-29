package widget

import (
	"github.com/duxphp/duxgo/core/ui/node"
)

type diy struct {
	node *node.TNode
}

// NewDiy 自定义节点
func NewDiy(node *node.TNode) *diy {
	return &diy{
		node: node,
	}
}

// SetData 设置数据
func (a *diy) SetData(data *map[string]any) {
	//a.data = data
}

// SetDialog 设置弹窗
func (a *diy) SetDialog(dialog bool) {
	//a.dialog = dialog
}

// Render 渲染UI
func (a *diy) Render() *node.TNode {
	return a.node
}
