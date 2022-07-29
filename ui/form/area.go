package form

import (
	"github.com/duxphp/duxgo/core/ui/node"
	"github.com/spf13/cast"
)

type AreaMap struct {
	Province string
	City     string
	Region   string
	Street   string
}

// Area 日期输入框
type Area struct {
	maps *AreaMap
}

// NewArea 创建日期
func NewArea(maps *AreaMap) *Area {
	if maps == nil {
		maps = &AreaMap{
			Province: "province",
			City:     "city",
			Region:   "region",
			Street:   "street",
		}
	}
	return &Area{
		maps: maps,
	}
}

// GetValue 格式化值
func (a *Area) GetValue(value any, info map[string]any) any {
	values := []any{}
	if a.maps.Province != "" {
		values = append(values, info[a.maps.Province])
	}
	if a.maps.City != "" {
		values = append(values, info[a.maps.City])
	}
	if a.maps.Region != "" {
		values = append(values, info[a.maps.Region])
	}
	if a.maps.Street != "" {
		values = append(values, info[a.maps.Street])
	}
	return values
}

// SaveValue 保存数据
func (a *Area) SaveValue(value any, data map[string]any) any {
	vals := cast.ToSlice(value)
	if a.maps.Province != "" {
		if len(vals) > 0 {
			data[a.maps.Province] = vals[0]
		} else {
			data[a.maps.Province] = ""
		}
	}
	if a.maps.City != "" {
		if len(vals) > 1 {
			data[a.maps.City] = vals[1]
		} else {
			data[a.maps.City] = ""
		}
	}
	if a.maps.Region != "" {
		if len(vals) > 2 {
			data[a.maps.Region] = vals[2]
		} else {
			data[a.maps.Region] = ""
		}
	}
	if a.maps.Street != "" {
		if len(vals) > 3 {
			data[a.maps.Street] = vals[3]
		} else {
			data[a.maps.Street] = ""
		}
	}
	return value
}

// Render 渲染
func (a *Area) Render(element node.IField) *node.TNode {
	ui := node.TNode{
		"nodeName": "app-cascader",
		"nParams": map[string]any{
			"allow-search": true,
			"path-mode":    true,
			"placeholder":  "请输入" + element.GetName(),
		},
		"dataUrl":      "/tools/area?level=3",
		"vModel:value": element.GetUIField(),
	}
	return &ui
}
