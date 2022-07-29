package table

import (
	"github.com/duxphp/duxgo/core/ui/node"
)

// IAction 动作接口
type IAction interface {
	Render() node.TNode
}

// Action 动作结构
type Action struct {
	UI IAction
}

// SetUI 设置链接
func (a *Action) SetUI(link IAction) *Action {
	a.UI = link
	return a
}

// Render 渲染UI
func (a *Action) Render() node.TNode {
	return a.UI.Render()
}
