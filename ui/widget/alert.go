package widget

import (
	"github.com/duxphp/duxgo/core/ui/node"
)

const (
	AlertInfo    Type = "info"
	AlertSuccess Type = "success"
	AlertWarning Type = "warning"
	AlertError   Type = "error"
)

type Type string

type alert struct {
	content string
	title   string
	mode    Type
}

func NewAlert(content string, title ...string) *alert {
	t := ""
	if len(title) > 0 {
		t = title[0]
	}
	return &alert{
		content: content,
		title:   t,
		mode:    "info",
	}
}

func (a *alert) SetType(name Type) *alert {
	a.mode = name
	return a
}

// Render 渲染UI
func (a *alert) Render() *node.TNode {
	return &node.TNode{
		"nodeName": "a-alert",
		"title":    a.title,
		"type":     a.mode,
		"child":    a.content,
	}
}
