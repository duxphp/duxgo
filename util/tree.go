package util

import (
	"fmt"
	"github.com/duxphp/duxgo/core"
	"github.com/duxphp/duxgo/response"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type TreeModel struct {
	Sort     uint `json:"sort" gorm:"default:0;notnull"`
	ParentID uint `json:"parent_id" gorm:"default:NULL"`
}

// TreeSort 树形排序控制器
func TreeSort(model any, key string) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		type Params struct {
			ID     int `json:"id"`
			Parent int `json:"parent"`
			Before int `json:"before"`
			After  int `json:"after"`
		}
		params := Params{}
		ctx.Bind(&params)

		tree := TreeSortT{
			model: model,
			key:   key,
		}

		fmt.Println(params, "params")

		var err error

		// 查询当前节点数据
		info := map[string]any{}
		core.Db.Model(model).Find(&info, params.ID)

		// 获取新位置上级ID
		var parentId int
		node := map[string]any{}
		if params.Before != 0 {
			err = core.Db.Model(model).Find(&node, params.Before).Error
			if err != nil {
				return err
			}
			parentId = cast.ToInt(node["parent_id"])
		} else if params.Parent != 0 {
			parentId = params.Parent
		}

		// 获取新位置ID列表
		ids := []int{}
		ids, err = tree.GetNodeIds(parentId)
		if err != nil {
			return err
		}

		// 从新位置删除当前ID
		if lo.IndexOf[int](ids, params.ID) != -1 {
			ids = DeleteSlice(ids, params.ID)
		}

		// 插入当前ID到指定位置
		var pos int
		for i, id := range ids {
			if params.Before == id {
				pos = i + 1
				break
			}
		}
		ids = append(ids, 0)
		copy(ids[pos+1:], ids[pos:])
		ids[pos] = params.ID

		// 新节点排序
		err = tree.UpdateNodeIds(ids, params.ID, parentId)
		if err != nil {
			return err
		}

		// 旧节点重排序
		if cast.ToInt(info["parent_id"]) != params.Parent {
			ids, err = tree.GetNodeIds(cast.ToInt(info["parent_id"]))
			if err != nil {
				return err
			}
			err = tree.UpdateNodeIds(ids, 0, 0)
			if err != nil {
				return err
			}
		}
		return response.New(ctx).Send("ok")
	}
}

type TreeSortT struct {
	model any
	key   string
}

// UpdateNodeIds 按ID更新数据
func (t *TreeSortT) UpdateNodeIds(ids []int, id int, parentId int) error {
	for i, v := range ids {
		data := map[string]any{
			"sort": i + 1,
		}
		if v == id {
			data["parent_id"] = lo.Ternary[any](parentId == 0, nil, parentId)
		}
		err := core.Db.Model(t.model).Where(t.key+" = ?", v).Updates(data).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// GetNodeIds 按ID获取数据
func (t *TreeSortT) GetNodeIds(parentId int) ([]int, error) {
	ids := []int{}
	modelDB := core.Db.Model(t.model)
	if parentId == 0 {
		modelDB.Where("parent_id is NULL")
	} else {
		modelDB.Where("parent_id = ?", parentId)
	}
	err := modelDB.Order("sort asc").Pluck(t.key, &ids).Error
	if err != nil {
		return ids, err
	}
	return ids, nil
}

// DeleteSlice 删除切片元素
func DeleteSlice(a []int, elem int) []int {
	j := 0
	for _, v := range a {
		if v != elem {
			a[j] = v
			j++
		}
	}
	return a[:j]
}
