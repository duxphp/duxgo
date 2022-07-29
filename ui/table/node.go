package table

import (
	"github.com/duxphp/duxgo/core/ui/node"
	"github.com/duxphp/duxgo/core/util/function"
)

type Side struct {
	Direction string
	Node      node.TNode
	Resize    bool
	Width     string
	Title     string
	Header    []node.TNode
}

// Node 表格结构
type Node struct {
	Url          string
	Primary      string
	DataFilter   map[string]any // 筛选数据
	Columns      []node.TNode   // 列数节点
	QuickFilters []node.TNode   // 快速筛选节点
	Filters      []node.TNode   // 筛选节点
	Tabs         []node.TNode   // tab节点
	Actions      []node.TNode   // 动作节点
	Bath         []node.TNode   // 批量操作节点
	Side         []*Side        // 动作节点
}

// SetUrl 设置数据链接
func (t *Node) SetUrl(url string) *Node {
	t.Url = url
	return t
}

// SetColumns 设置列
func (t *Node) SetColumns(node []node.TNode) *Node {
	t.Columns = node
	return t
}

// SetDataFilter 设置筛选
func (t *Node) SetDataFilter(data map[string]any) *Node {
	t.DataFilter = data
	return t
}

// SetFilters 设置筛选
func (t *Node) SetFilters(node []node.TNode) *Node {
	t.Filters = node
	return t
}

// SetQuickFilters 设置快速筛选
func (t *Node) SetQuickFilters(node []node.TNode) *Node {
	t.QuickFilters = node
	return t
}

// SetTabs 设置切换筛选
func (t *Node) SetTabs(node []node.TNode) *Node {
	t.Tabs = node
	return t
}

// SetActions 设置动作
func (t *Node) SetActions(node []node.TNode) *Node {
	t.Actions = node
	return t
}

// SetSide 设置侧边栏
func (t *Node) SetSide(node []*Side) *Node {
	t.Side = node
	return t
}

// HeaderNode 表格头数据
func (t *Node) HeaderNode() []node.TNode {

	rightNode := t.QuickFilters
	if len(t.Filters) > 0 {
		filterExpand := node.TNode{
			"nodeName": "a-trigger",
			"trigger":  "click",
			"child": []node.TNode{
				{
					"nodeName": "a-button",
					"type":     "secondary",
					"child": []node.TNode{
						{
							"nodeName": "span",
							"child":    "筛选",
						},
						{
							"vSlot:icon": "",
							"nodeName":   "icon-filter",
						},
					},
				},
				{
					"vSlot:content": "",
					"nodeName":      "div",
					"class":         "flex flex-col rounded shadow bg-white dark:bg-blackgray-1 dark:text-gray-400 p-2 w-56",
					"child":         t.Filters,
				},
			},
		}
		rightNode = append(rightNode, filterExpand)
	}

	var tabsNode node.TNode
	if len(t.Tabs) > 0 {
		tabsNode = node.TNode{
			"nodeName":          "a-radio-group",
			"name":              "type",
			"type":              "button",
			"vModel:modelValue": "data.filter.type",
			"child":             t.Tabs,
		}
	}

	for _, tNode := range t.Actions {
		rightNode = append(rightNode, tNode)
	}

	return []node.TNode{
		{
			"nodeName": "div",
			"class":    "flex-grow w-10 flex justify-start",
			"child":    tabsNode,
		},
		{
			"nodeName": "div",
			"class":    "flex-none flex gap-2",
			"child":    rightNode,
		},
	}
}

// TableNode 表格节点
func (t *Node) TableNode() *node.TNode {

	return &node.TNode{
		"nodeName":         "app-table",
		"requestEventName": function.Md5(t.Url),
		"class":            "",
		"url":              t.Url,
		"urlBind":          true,
		"nowrap":           true,
		"nParams": node.TNode{
			"row-key": t.Primary,
			//"scroll": node.TNode{
			//	"x": "110%",
			//},
		},
		"columns":         &t.Columns,
		"vBind:filter":    "data.filter",
		"select":          t.Bath != nil && len(t.Bath) > 0,
		"tablLayoutFixed": true,
		"child": node.TNode{
			"vSlot:footer": "footer",
			"nodeName":     "div",
			"class":        "flex gap-2",
			"child":        t.Bath,
		},
	}
}

// Render 树形渲染
func (t *Node) Render() *node.TNode {
	var bodyNode []node.TNode
	for _, side := range t.Side {
		if side.Direction != "left" {
			continue
		}
		var sideNode node.TNode
		nodeName := ""
		style := ""
		directions := ""
		class := "border-r border-gray-200 dark:border-gray-700 flex-none bg-white dark:bg-blackgray-4 h-screen p-2 flex flex-col"
		style = "width:" + side.Width
		if side.Resize {
			nodeName = "a-resize-box"
			directions = "right"
		} else {
			nodeName = "div"
		}

		sideChild := []node.TNode{}
		// 侧栏标题
		if side.Title != "" {
			sideChild = append(sideChild, node.TNode{
				"nodeName": "div",
				"class":    "py-2 text-gray-500",
				"child":    side.Title,
			})
		}

		if len(side.Header) > 0 {
			sideChild = append(sideChild, side.Header...)
		}
		sideChild = append(sideChild, node.TNode{
			"nodeName": "div",
			"class":    "flex-grow h-10 mt-2",
			"child":    side.Node,
		})

		sideNode = node.TNode{
			"nodeName":   nodeName,
			"style":      style,
			"directions": []string{directions},
			"class":      class,
			"child":      sideChild,
		}
		bodyNode = append(bodyNode, sideNode)
	}

	bodyNode = append(bodyNode, node.TNode{
		"nodeName": "app-layout",
		"class":    "flex-grow w-10",
		"title":    "数据列表",
		"child": node.TNode{
			"vSlot":    "",
			"nodeName": "div",
			"class":    "flex flex-row items-start gap-4 p-4",
			"child": []node.TNode{
				{
					"nodeName": "div",
					"class":    "flex-grow lg:w-10 p-4 bg-white dark:bg-blackgray-4 rounded shadow",
					"child": []node.TNode{
						{
							"nodeName": "div",
							"class":    "flex-none flex flex-row gap-2 items-center pb-4",
							"child":    t.HeaderNode(),
						},
						*t.TableNode(),
					},
				},
			},
		},
	})

	for _, side := range t.Side {
		if side.Direction != "right" {
			continue
		}
		var sideNode node.TNode
		nodeName := ""
		style := ""
		directions := ""
		class := "border-r border-gray-200 dark:border-gray-700 flex-none bg-white dark:bg-blackgray-4 h-screen"
		style = "width:" + side.Width
		if side.Resize {
			nodeName = "a-resize-box"
			directions = "left"
		} else {
			nodeName = "div"
		}
		sideNode = node.TNode{
			"nodeName":   nodeName,
			"style":      style,
			"directions": []string{directions},
			"class":      class,
			"child":      side.Node,
		}
		bodyNode = append(bodyNode, sideNode)
	}

	return &node.TNode{
		"node": node.TNode{
			"nodeName": "app-form",
			"value": map[string]any{
				"filter": t.DataFilter,
			},
			"child": node.TNode{
				"nodeName": "div",
				"class":    "flex h-screen",
				"vSlot":    "{value: data}",
				"child":    bodyNode,
			},
		},
		"setupScript": "",
	}
}
