package widget

import (
	"github.com/duxphp/duxgo/ui/node"
	"github.com/duxphp/duxgo/util/function"
)

type Link struct {
	name         string
	url          string
	params       map[string]any
	icon         string
	mode         string
	modeConfig   node.TNode
	buttonType   string
	buttonStatus string
	buttonLong   bool
	prefix       string
	fields       []string
}

// NewLink 新建链接
func NewLink(name string, url string, params ...map[string]any) *Link {
	link := Link{
		name: name,
		url:  url,
	}
	if len(params) > 0 {
		link.params = params[0]
	} else {
		link.params = map[string]any{}
	}
	return &link
}

// SetType 设置类型
func (a *Link) SetType(name string, config ...node.TNode) *Link {
	a.mode = name
	if len(config) > 0 {
		a.modeConfig = config[0]
	}
	return a
}

// SetIcon 设置图标
func (a *Link) SetIcon(icon string) *Link {
	a.icon = icon
	return a
}

// SetButton 设置按钮
func (a *Link) SetButton(params ...any) *Link {
	if len(params) > 0 {
		a.buttonType = params[0].(string)
	} else {
		a.buttonType = "primary"
	}
	if len(params) > 1 {
		a.buttonStatus = params[1].(string)
	} else {
		a.buttonStatus = "normal"
	}
	if len(params) > 2 {
		a.buttonLong = params[2].(bool)
	} else {
		a.buttonLong = false
	}
	return a
}

// SetModel 设置模型链接
func (a *Link) SetModel(prefix string, fields []string) *Link {
	a.prefix = prefix
	a.fields = fields
	return a
}

// Render 渲染UI
func (a *Link) Render() node.TNode {
	url := function.BuildUrl(a.url, a.params, false, a.prefix, a.fields)

	nodeData := node.TNode{
		"nodeName": "route",
	}

	//链接类型
	switch a.mode {
	case "blank":
		nodeData["nodeName"] = "a"
		nodeData["vBind:href"] = url
		nodeData["target"] = "_blank"
	case "dialog":
		nodeData["vBind:href"] = url
		nodeData["type"] = "dialog"
		nodeData["title"] = a.name
	case "drawer":
		nodeData["vBind:href"] = url
		nodeData["type"] = "dialog"
		nodeData["mode"] = "drawer"
		nodeData["title"] = a.name
	case "ajax":
		nodeData["vBind:href"] = url
		nodeData["type"] = "ajax"
		nodeData["title"] = "确认进行" + a.name + "操作？"
	default:
		nodeData["vBind:href"] = url
	}

	for key, value := range a.modeConfig {
		nodeData[key] = value
	}

	// 链接内部元素
	var childData []node.TNode
	childData = append(childData, node.TNode{
		"nodeName": "span",
		"child":    a.name,
	})
	var linkData node.TNode
	if a.buttonType != "" {
		// 按钮类型
		if a.icon != "" {
			childData = append(childData, node.TNode{
				"vSlot:icon": "",
				"nodeName":   "icon-" + a.icon,
			})
		}
		linkData = node.TNode{
			"nodeName": "a-button",
			"type":     a.buttonType,
			"status":   a.buttonStatus,
			"long":     a.buttonLong,
			"child":    childData,
		}
	} else {
		// 普通类型
		if a.icon != "" {
			childData = append(childData, node.TNode{
				"class":    "mr-2",
				"nodeName": "icon-" + a.icon,
			})
		}
		linkData = node.TNode{
			"nodeName": "span",
			"class":    "arco-link arco-link-status-normal",
			"child":    childData,
		}
	}
	nodeData["child"] = linkData

	return nodeData
}
