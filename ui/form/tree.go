package form

import (
	"github.com/duxphp/duxgo/core/ui/node"
	"encoding/json"
)

// Tree 树形选择框
type Tree struct {
	num  int
	data []map[string]any
}

// StructType 结构体类型
func (a *Tree) StructType() any {
	return ""
}

// NewTree 创建树形
func NewTree() *Tree {
	return &Tree{}
}

// SetData 设置数据
func (a *Tree) SetData(data []map[string]any) *Tree {
	a.data = data
	return a
}

// GetValue 格式化值
func (a *Tree) GetValue(value any, info map[string]any) any {
	return value
}

// SaveValue 保存数据
func (a *Tree) SaveValue(value any, data map[string]any) any {
	marshal, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	return marshal
}

// Render 渲染
func (a *Tree) Render(element node.IField) *node.TNode {
	ui := node.TNode{
		"nodeName": "div",
		"class":    "bg-gray-100 dark:bg-blackgray-2 p-2 rounded w-full h-56  overflow-y-auto app-scrollbar",
		"child": map[string]any{
			"nodeName":            "a-tree",
			"blockNode":           true,
			"checkable":           true,
			"showLine":            true,
			"vModel:checked-keys": element.GetUIField(),
			"data":                treeLoop(a.data),
		},
	}
	return &ui
}

func treeLoop(data []map[string]any) *[]map[string]any {
	var dataArr []map[string]any
	for _, item := range data {
		tmpData := map[string]any{
			"title": item["name"],
			"key":   item["id"],
		}
		if item["children"] != nil {
			tmpData["children"] = treeLoop(item["children"].([]map[string]any))
		}
		dataArr = append(dataArr, tmpData)
	}
	return &dataArr
}
