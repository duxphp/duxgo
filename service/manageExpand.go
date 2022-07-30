package service

import (
	"github.com/duxphp/duxgo/core"
	"github.com/duxphp/duxgo/exception"
	"github.com/duxphp/duxgo/response"
	"github.com/duxphp/duxgo/ui/form"
	"github.com/duxphp/duxgo/ui/table"
	function2 "github.com/duxphp/duxgo/util/function"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type ManageExpand struct {
	event       string
	table       func(echo.Context) *table.Table
	form        func(echo.Context) *form.Form
	delCall     func(id string, model *gorm.DB) error
	model       any
	modelKey    string
	searchWhere func(ctx echo.Context, model *gorm.DB) (*gorm.DB, error)
	searchField []string
	searchMaps  map[string]any
}

// NewManageExpand 实例化扩展方法
func NewManageExpand(event ...string) *ManageExpand {
	var eventName string
	if len(event) > 0 {
		eventName = event[0]
	}
	return &ManageExpand{
		event: function2.Md5(eventName),
	}
}

// SetTable 设置表格
func (t *ManageExpand) SetTable(table func(echo.Context) *table.Table) *ManageExpand {
	t.table = table
	return t
}

// SetForm 设置表单
func (t *ManageExpand) SetForm(form func(echo.Context) *form.Form) *ManageExpand {
	t.form = form
	return t
}

// SetModel 设置模型
func (t *ManageExpand) SetModel(model any, key string) *ManageExpand {
	t.model = model
	t.modelKey = key
	return t
}

// CallDel 设置删除毁掉
func (t *ManageExpand) CallDel(callback func(id string, model *gorm.DB) error) *ManageExpand {
	t.delCall = callback
	return t
}

// ListPage 列表页面
func (t *ManageExpand) ListPage(ctx echo.Context) error {
	return response.New(ctx).Send("ok", t.table(ctx).Render(ctx))
}

// ListData 列表数据
func (t *ManageExpand) ListData(ctx echo.Context) error {
	return response.New(ctx).Send("ok", t.table(ctx).Data(ctx))
}

// FormPage 表单页面
func (t *ManageExpand) FormPage(ctx echo.Context) error {
	return response.New(ctx).Send("ok", t.form(ctx).Render(ctx))
}

// FormSave 表单保存
func (t *ManageExpand) FormSave(ctx echo.Context, eventPos ...string) error {
	formUI := t.form(ctx)
	key := formUI.GetKey()
	err := formUI.Save(ctx)

	if err != nil {
		return err
	}
	mode := "add"
	if key > 0 {
		mode = "edit"
	}
	event := map[string]any{}

	if t.table == nil {
		return response.New(ctx).Send("保存数据成功", event)
	}

	// 插入位置
	pos := ""
	if len(eventPos) > 0 {
		pos = eventPos[0]
	}

	// 获取表格数据
	tableUI := t.table(ctx)
	data := tableUI.Data(ctx)

	// 设置父级id
	var parentKey uint
	if tableUI.GetTree() {
		info := formUI.GetInfo()
		parentKey = cast.ToUint(info["parent_id"])
	}
	// 获取节点数据
	tree := function2.GetTreeNode(data["data"], formUI.GetKey(), "id", "children")
	data["data"] = []map[string]any{tree}

	eventData := []map[string]any{}
	for _, item := range data["data"].([]map[string]any) {
		eventData = append(eventData, map[string]any{
			"type":      mode,
			"key":       item[formUI.GetPrimary()],
			"data":      item,
			"pos":       pos,
			"parentKey": parentKey,
		})
	}
	event = map[string]any{
		"__event": map[string]any{
			"name": t.event,
			"data": eventData,
		},
	}
	return response.New(ctx).Send("保存数据成功", event)
}

// Status 状态更改
func (t *ManageExpand) Status(ctx echo.Context, model any) error {
	id := ctx.QueryParam("id")
	if id == "" {
		return exception.BusinessError("参数传递错误")
	}
	body := function2.CtxBody(ctx)
	field := gjson.GetBytes(body, "field").String()
	value := gjson.GetBytes(body, "status").Bool()

	core.Db.First(model, id)
	core.Db.Model(model).Update(field, value)
	return response.New(ctx).Send("更改状态成功")
}

// Del 删除数据
func (t *ManageExpand) Del(ctx echo.Context, model any) error {
	id := ctx.QueryParam("id")
	if id == "" {
		return exception.BusinessError("参数传递错误")
	}
	event := map[string]any{
		"__event": map[string]any{
			"name": t.event,
			"data": []map[string]any{
				{
					"type": "del",
					"key":  cast.ToUint(id),
					"data": []map[string]any{},
				},
			},
		},
	}
	tx := core.Db.Begin()
	if t.delCall != nil {
		err := t.delCall(id, tx)
		if err != nil {
			return err
		}
	}
	tx.Delete(model, id)
	tx.Commit()
	return response.New(ctx).Send("删除数据成功", event)
}

type SelectParams struct {
	Id    string `query:"id"`
	Query string `query:"query"`
}

type SearchMapFun func(map[string]any) string

func (t *ManageExpand) SearchMaps(fields map[string]any) *ManageExpand {
	t.searchMaps = fields
	return t
}

func (t *ManageExpand) SearchField(fields []string) *ManageExpand {
	t.searchField = fields
	return t
}

func (t *ManageExpand) SearchWhere(call func(ctx echo.Context, model *gorm.DB) (*gorm.DB, error)) *ManageExpand {
	t.searchWhere = call
	return t
}

func (t *ManageExpand) Search(ctx echo.Context) error {

	params := SelectParams{}
	err := ctx.Bind(&params)
	if err != nil {
		return err
	}

	model := core.Db.Model(t.model).Debug()

	if params.Query != "" && len(t.searchField) > 0 {
		where := core.Db
		for _, field := range t.searchField {
			where = where.Or(field+` like ?`, "%"+params.Query+"%")
		}
		model = model.Where(where.Or(t.modelKey+` = ?`, params.Query))
	}

	if params.Id != "" {
		ids := strings.Split(params.Id, ",")
		model = model.Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "FIELD(" + t.modelKey + ",?)", Vars: []any{ids}, WithoutParentheses: true},
		})
	}

	if t.searchWhere != nil {
		model, err = t.searchWhere(ctx, model)
		if err != nil {
			return err
		}
	}

	var count int64
	err = model.Count(&count).Error

	type page struct {
		Page  int `query:"page"`
		Limit int `query:"limit"`
	}
	pageQuery := page{
		Page:  1,
		Limit: 50,
	}
	_ = ctx.Bind(&pageQuery)
	offset, totalPage := function2.PageLimit(pageQuery.Page, cast.ToInt(count), pageQuery.Limit)

	if err != nil {
		return err
	}
	datas := []map[string]any{}
	model.Limit(pageQuery.Limit).Offset(offset).Find(&datas)

	if t.searchMaps == nil {
		t.searchMaps = map[string]any{"name": "name", "id": "id"}
	}

	result := []map[string]any{}
	for _, data := range datas {
		itemData := map[string]any{}
		for k, v := range t.searchMaps {
			switch v.(type) {
			case string:
				itemData[k] = data[v.(string)]
			case SearchMapFun:
				itemData[k] = v.(SearchMapFun)(data)
			}
		}
		result = append(result, itemData)
	}

	return response.New(ctx).Send("ok", map[string]any{
		"data":      result,
		"pageSize":  pageQuery.Limit,
		"total":     count,
		"totalPage": totalPage,
	})
}
