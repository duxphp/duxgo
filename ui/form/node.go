package form

import (
	"github.com/duxphp/duxgo/ui/node"
)

// Node 表单结构
type Node struct {
	Url     string
	Method  string
	Title   string
	Data    *map[string]any
	Element []*node.TNode
	Header  []node.IWidget
	Back    bool
	Dialog  bool
}

// SetBack 设置返回元素
func (t *Node) SetBack(back bool) *Node {
	t.Back = back
	return t
}

// SetDialog 设置表格元素
func (t *Node) SetDialog(dialog bool) *Node {
	t.Dialog = dialog
	return t
}

// SetHeader 设置头部元素
func (t *Node) SetHeader(widgets []node.IWidget) *Node {
	t.Header = widgets
	return t
}

// SetElement 设置表格元素
func (t *Node) SetElement(elm []*node.TNode) *Node {
	t.Element = elm
	return t
}

// SetData 设置表单数据
func (t *Node) SetData(data *map[string]any) *Node {
	t.Data = data
	return t
}

// Render 表单渲染
func (t *Node) Render() *node.TNode {
	var element []*node.TNode
	if len(t.Header) > 0 {
		for _, widget := range t.Header {
			element = append(element, &node.TNode{
				"nodeName": "div",
				"class":    "mb-4",
				"child":    widget.Render(),
			})
		}
	}
	t.Element = append(element, t.Element...)
	var renderForm any

	if t.Dialog {
		renderForm = t.renderDialog()
	} else {
		renderForm = t.renderPage()
	}

	return &node.TNode{
		"node": node.TNode{
			"nodeName": "app-form",
			"url":      t.Url,
			"method":   t.Method,
			"value":    t.Data,
			"layout":   "vertical",
			"back":     t.Back,
			"child": node.TNode{
				"nodeName": "div",
				"class":    "flex",
				"vSlot":    "{value: data, submitStatus: loading}",
				"child":    renderForm,
			},
		},
		"setupScript": "",
	}
}

// renderPage 渲染页面表单
func (t *Node) renderPage() *[]node.TNode {
	backNode := node.TNode{}
	submitText := "保存"
	if t.Back {
		backNode = node.TNode{
			"nodeName": "route",
			"type":     "back",
			"child": node.TNode{
				"type":     "outline",
				"nodeName": "a-button",
				"child":    "返回",
			},
		}
		submitText = "提交"
	}

	return &[]node.TNode{
		// 边栏
		// 主体
		{
			"nodeName":          "app-layout",
			"class":             "flex-grow w-10",
			"title":             t.Title,
			"form":              true,
			"back":              t.Back,
			"vBind:formLoading": "loading",
			"child": []node.TNode{
				{
					"nodeName": "div",
					"class":    "p-4",
					"child": []node.TNode{
						{
							"nodeName": "div",
							"child":    t.Element,
						},
						{
							"nodeName": "div",
							"class":    "flex items-center justify-end gap-2 flex-row",
							"child": []node.TNode{
								backNode,
								{
									"nodeName":      "a-button",
									"html-type":     "submit",
									"vBind:loading": "loading",
									"type":          "primary",
									"child":         submitText,
								},
							},
						},
					},
				},
			},
		},
		// 边栏
	}
}

// renderPage 渲染弹窗表单
func (t *Node) renderDialog() *node.TNode {
	return &node.TNode{
		"nodeName": "app-dialog",
		"title":    t.Title,
		"class":    "flex-grow",
		"child": []node.TNode{
			{
				"nodeName":      "div",
				"vSlot:default": "",
				"class":         "flex",
				"child": []node.TNode{
					{
						"nodeName": "div",
						"class":    "flex-grow p-5 pb-0",
						"child":    t.Element,
					},
				},
			},
			{
				"nodeName":     "div",
				"vSlot:footer": "",
				"class":        "arco-modal-footer",
				"child": []node.TNode{
					{
						"nodeName": "route",
						"type":     "back",
						"child": node.TNode{
							"nodeName": "a-button",
							"child":    "取消",
						},
					},
					{
						"nodeName":      "a-button",
						"type":          "primary",
						"html-type":     "submit",
						"vBind:loading": "loading",
						"child":         "提交",
					},
				},
			},
		},
	}
}
