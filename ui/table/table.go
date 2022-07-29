package table

import (
	"encoding/json"
	"fmt"
	"github.com/duxphp/duxgo/global"
	"github.com/duxphp/duxgo/ui/node"
	function2 "github.com/duxphp/duxgo/util/function"
	"github.com/gookit/goutil/maputil"
	"github.com/jianfengye/collection"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type collectOrder struct {
	key  string
	desc bool
}

type Table struct {
	model        any
	modelDB      *gorm.DB
	collect      collection.ICollection
	collectFun   func(filter map[string]any) collection.ICollection
	collectOrder []collectOrder
	modelOrder   []string
	url          string
	fieldMaps    map[string]string
	back         bool
	primary      string
	actions      []*Action
	columns      []*Column
	filters      []*Filter
	tabs         []*Tab
	sides        []*Side
	tree         bool
	limit        int
}

// NewTable 实例化表格UI
func NewTable() *Table {
	return &Table{}
}

// SetUrl 设置数据链接
func (t *Table) SetUrl(url string) *Table {
	t.url = url
	return t
}

// AddFields 添加查询字段
func (t *Table) AddFields(data map[string]string) *Table {
	t.fieldMaps = data
	return t
}

// SetData 设置数据模式
func (t *Table) SetData(collect collection.ICollection, primary string) *Table {
	t.collect = collect
	t.primary = primary
	return t
}

// SetDataFun 设置数据回调模式
func (t *Table) SetDataFun(collectFun func(filter map[string]any) collection.ICollection, primary string) *Table {
	t.collectFun = collectFun
	t.primary = primary
	return t
}

// DataOrder 设置排序规则
func (t *Table) DataOrder(key string, desc ...bool) *Table {
	sort := false
	if len(desc) > 0 {
		sort = desc[0]
	}
	t.collectOrder = append(t.collectOrder, collectOrder{key: key, desc: sort})
	return t
}

// SetModel 设置模型模式
func (t *Table) SetModel(mode any, primary string) *Table {
	t.model = mode
	t.modelDB = global.Db.Model(mode)
	t.primary = primary
	return t
}

func (t *Table) SetLimit(limit int) *Table {
	t.limit = limit
	return t
}

// SetTree 设置树形数据
func (t *Table) SetTree() *Table {
	t.tree = true
	return t
}

// GetTree 获取树形状态
func (t *Table) GetTree() bool {
	return t.tree
}

// ModelOrder 设置排序规则
func (t *Table) ModelOrder(order string) *Table {
	t.modelOrder = append(t.modelOrder, order)
	return t
}

// PreloadModel 预载模型
func (t *Table) PreloadModel(query string, args ...any) *Table {
	t.modelDB.Preload(query, args)
	return t
}

// GetModel 获取模型
func (t *Table) GetModel() *gorm.DB {
	return t.modelDB
}

// GetRowField 获取行字段
func (t *Table) GetRowField(field string) string {
	return "rowData." + field
}

// AddCol 添加表格列
func (t *Table) AddCol(name string, field string, callback ...func(val any, items map[string]any) any) *Column {
	elm := Column{
		Name:  name,
		Field: field,
	}
	if len(callback) > 0 {
		elm.Format = callback[0]
	}
	t.columns = append(t.columns, &elm)
	return &elm
}

// AddFilter 添加筛选
func (t *Table) AddFilter(name string, field string) *Filter {
	elm := Filter{
		Name:    name,
		Field:   field,
		Quick:   false,
		Default: nil,
	}
	t.filters = append(t.filters, &elm)
	return &elm
}

// AddTab 添加切换类型
func (t *Table) AddTab(name string, where ...func(*gorm.DB)) *Tab {
	elm := Tab{
		Name: name,
	}
	if where != nil {
		elm.Where = where[0]
	}
	t.tabs = append(t.tabs, &elm)
	return &elm
}

// AddAction 添加操作
func (t *Table) AddAction() *Action {
	elm := Action{}
	t.actions = append(t.actions, &elm)
	return &elm
}

func (t *Table) AddSide(side *Side) *Table {
	t.sides = append(t.sides, side)
	return t
}

// Render 渲染表格
func (t *Table) Render(ctx echo.Context) *node.TNode {

	fields := []string{
		t.primary,
	}
	for _, item := range t.columns {
		fields = append(fields, item.Field)
	}

	var cols []node.TNode
	for _, item := range t.columns {
		cols = append(cols, item.Render(fields))
	}

	// 筛选提取
	var filters []node.TNode
	var QuickFilters []node.TNode
	DataFilter := map[string]any{}
	for _, item := range t.filters {
		var value any
		value = ctx.QueryParam(item.Field)
		if value == "" {
			value = item.Default
		}
		DataFilter[item.Field] = value
		if item.UI == nil {
			continue
		}
		if item.Quick {
			QuickFilters = append(QuickFilters, item.Render())
		} else {
			filters = append(filters, item.Render())
		}
	}

	// Tab提取
	var tabs []node.TNode
	for index, item := range t.tabs {
		tabs = append(tabs, item.Render(index))
	}
	if len(t.tabs) > 0 {
		DataFilter["type"] = 0
	}

	// 动作提取
	var actions []node.TNode
	for _, item := range t.actions {
		actions = append(actions, item.Render())
	}

	tableNode := &Node{
		Url:     t.url,
		Primary: t.primary,
	}

	var side []*Side
	for _, item := range t.sides {
		side = append(side, item)
	}

	tableNode.SetQuickFilters(QuickFilters)
	tableNode.SetFilters(filters)
	tableNode.SetTabs(tabs)
	tableNode.SetActions(actions)
	tableNode.SetColumns(cols)
	tableNode.SetDataFilter(DataFilter)
	tableNode.SetSide(side)

	return tableNode.Render()
}

// Data 渲染数据
func (t *Table) Data(ctx echo.Context) map[string]any {
	if t.modelDB == nil && t.collect == nil && t.collectFun == nil {
		panic("table Model or Data not set")
	}

	model := t.modelDB
	collect := t.collect
	collectFilter := map[string]any{}

	// 条件处理
	where := map[string]any{}
	for _, filter := range t.filters {
		value := ctx.QueryParam(filter.Field)
		if value != "" {
			// 集合条件
			collectFilter[filter.Field] = value
			// 模型处理
			if t.modelDB != nil {
				if filter.Where == nil {
					where[filter.Field] = value
				} else {
					filter.Where(value, t.modelDB)
				}
			}
			// 集合处理
			if t.collect != nil {
				if filter.Collect == nil {
					filter.Collect(&collect)
				}
			}
		}
	}

	if len(where) > 0 && t.modelDB != nil {
		model.Where(where)
	}
	// tab处理
	tabValue, _ := strconv.Atoi(ctx.QueryParam("type"))
	for index, tab := range t.tabs {
		if tabValue == index {
			collectFilter["tab"] = index
			if tab.Where != nil {
				// 模型处理
				if t.model != nil {
					tab.Where(t.modelDB)
				}
				// 集合处理
				if t.collect != nil {
					tab.Collect(&collect)
				}
			}
		}
	}
	// 排序处理
	for _, column := range t.columns {
		if !column.Sort || column.Field == "" {
			continue
		}
		getSort := ctx.QueryParam(fmt.Sprintf("_sort[%s]", column.Field))
		if getSort == "" {
			continue
		}

		if t.model != nil {
			model.Order(column.Field + " " + getSort)
		}
	}

	// 集合回调数据
	if t.collectFun != nil {
		collect = t.collectFun(collectFilter)
	}

	// 字段处理
	fields := []string{
		t.primary,
	}
	if t.tree {
		fields = append(fields, "parent_id")
	}
	formats := map[string]func(any, map[string]any) any{}
	for _, item := range t.columns {
		if item.Field == "" {
			continue
		}
		fields = append(fields, item.Field)
		if item.Format != nil {
			formats[item.Field] = item.Format
		}
	}

	// 分页处理
	var total int64
	if t.model != nil {
		model.Count(&total)
	}
	if t.collect != nil || t.collectFun != nil {
		total = int64(collect.Count())
	}
	if t.tree {
		t.limit = 10000
	}

	type page struct {
		Page  int `query:"page"`
		Limit int `query:"limit"`
	}
	pageQuery := page{
		Page:  1,
		Limit: t.limit,
	}
	_ = ctx.Bind(&pageQuery)
	offset, totalPage := function2.PageLimit(pageQuery.Page, cast.ToInt(total), pageQuery.Limit)

	// 数据查询
	var data []map[string]any
	if t.model != nil {
		if len(t.modelOrder) > 0 {
			for _, order := range t.modelOrder {
				model.Order(order)
			}
		}

		if t.tree {
			model.Order("sort asc")
		}

		err := model.Limit(pageQuery.Limit).Offset(offset).Find(&data).Error
		if err != nil {
			panic(err)
		}

		if t.tree && len(data) > 0 {
			data = function2.SliceToTree(data, t.primary, "parent_id", "children")
		}
		data = t.filterData(data, fields, formats)

	}
	if t.collect != nil || t.collectFun != nil {
		if t.collectOrder != nil {
			for _, order := range t.collectOrder {
				if order.desc {
					collect = collect.SortByDesc(order.key)
				} else {
					collect = collect.SortBy(order.key)
				}
			}
		}
		ret, _ := collect.Slice(offset, pageQuery.Limit).ToJson()
		json.Unmarshal(ret, &data)
		if t.tree && len(data) > 0 {
			data = function2.SliceToTree(data, t.primary, "parent_id", "children")
		}
		data = t.filterData(data, fields, formats)
	}

	return map[string]any{
		"data":      data,
		"pageSize":  pageQuery.Limit,
		"total":     total,
		"totalPage": totalPage,
	}
}

// filterData 过滤返回数据
func (t *Table) filterData(listData []map[string]any, fields []string, formats map[string]func(any, map[string]any) any) []map[string]any {

	// 过滤列表字段
	data := []map[string]any{}
	for _, items := range listData {
		newItems := map[string]any{}
		// 提取设置字段
		for key, val := range t.fieldMaps {
			newField := strings.Replace(key, ".", "_", -1)
			v, o := maputil.GetByPath(val, items)
			if !o {
				newItems[newField] = nil
			} else {
				newItems[newField] = v
			}
		}
		// 提取级联字段
		for _, field := range fields {
			newField := strings.Replace(field, ".", "_", -1)
			val, ok := maputil.GetByPath(field, items)
			if !ok {
				newItems[newField] = nil
			} else {
				newItems[newField] = val
			}
			if formats[field] != nil {
				newItems[newField] = formats[field](val, items)
			}
		}

		// 转换树形字段

		if _, ok := items["children"]; ok && t.tree {
			newItems["children"] = t.filterData(items["children"].([]map[string]any), fields, formats)
		}

		data = append(data, newItems)
	}
	return data
}

func TreePreload(d *gorm.DB) *gorm.DB {
	return d.Preload("Children", TreePreload)
}
