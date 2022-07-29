package form

import (
	"encoding/json"
	"fmt"
	exception2 "github.com/duxphp/duxgo/exception"
	"github.com/duxphp/duxgo/global"
	"github.com/duxphp/duxgo/ui/node"
	"github.com/duxphp/duxgo/util/function"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"reflect"
)

type Form struct {
	model      any
	info       map[string]any                                            // 表单数据
	dialog     bool                                                      // 弹出类型
	url        string                                                    // 提交URL
	method     string                                                    //提交方式
	title      string                                                    // 表单标题
	back       bool                                                      // 返回功能
	header     []node.IWidget                                            //头部部件
	element    []*Element                                                // 表单元素集
	modelDB    *gorm.DB                                                  // 模型结构体
	primary    string                                                    // 主键名
	key        uint                                                      // 主键id
	validate   *validator.Validate                                       //验证规则
	saveFn     func(data map[string]any, key uint) error                 // 保存函数
	saveBefore func(data map[string]any, update bool, db *gorm.DB) error // 保存函数
	saveAfter  func(model any, update bool) error                        // 保存函数

}

// NewForm 实例化表单UI
func NewForm() *Form {
	return &Form{
		info:     map[string]any{},
		method:   "post",
		title:    "信息详情",
		back:     true,
		dialog:   true,
		validate: global.Validator,
	}
}

// SetModel 设置模型
func (t *Form) SetModel(mode any, primary string, id ...uint) *Form {
	t.model = mode
	t.modelDB = global.Db.Model(t.model)
	t.primary = primary
	if len(id) > 0 {
		t.key = id[0]
	}
	return t
}

// SetData 设置表单数据
func (t *Form) SetData(data map[string]any) *Form {
	t.info = data
	return t
}

// SetDialog 设置弹窗
func (t *Form) SetDialog(status bool) *Form {
	t.dialog = status
	return t
}

// SetUrl 设置提交链接
func (t *Form) SetUrl(url string, method ...string) *Form {
	t.url = url
	if len(method) > 0 {
		t.method = method[0]
	}
	return t
}

// SetTitle 设置表单标题
func (t *Form) SetTitle(title string) *Form {
	t.title = title
	return t
}

// SetBack 设置返回
func (t *Form) SetBack(status bool) *Form {
	t.back = status
	return t
}

// RegValidator 注册验证规则
func (t *Form) RegValidator(tag string, fn validator.Func, callValidationEvenIfNull ...bool) *Form {
	err := t.validate.RegisterValidation(tag, fn, callValidationEvenIfNull...)
	if err != nil {
		panic(err.Error())
	}
	return t
}

// GetKey 获取主键
func (t *Form) GetKey() uint {
	return t.key
}

// GetPrimary 获取主键
func (t *Form) GetPrimary() string {
	return t.primary
}

// GetModel 获取模型
func (t *Form) GetModel() *gorm.DB {
	return t.modelDB
}

// GetInfo 获取模型数据
func (t *Form) GetInfo() map[string]any {
	return t.info
}

// AddField 添加元素
func (t *Form) AddField(name string, field string) *Element {
	elm := Element{
		Name:  name,
		Field: field,
	}
	t.element = append(t.element, &elm)
	return &elm
}

// AddHeader 添加头部元素
func (t *Form) AddHeader(widget node.IWidget) {
	t.header = append(t.header, widget)
}

// AddLayout 添加布局元素
func (t *Form) AddLayout(element ILayout, callback func(form *Form)) {
	element.SetData(t.info)
	element.SetDialog(t.dialog)
	elm := Element{
		Layout: element,
	}
	t.element = append(t.element, &elm)
	element.Column(callback)
}

// AddColumn 添加多列布局
func (t *Form) AddColumn(element ILayout, callback func(element ILayout)) {
	element.SetData(t.info)
	element.SetDialog(t.dialog)
	elm := Element{
		Layout: element,
	}
	t.element = append(t.element, &elm)
	callback(element.(ILayout))
}

// RenderElement 渲染列表
func (t *Form) RenderElement() []*node.TNode {
	var element []*node.TNode
	for _, item := range t.element {
		if item.UI != nil {
			// 普通UI
			element = append(element, item.Render())
		}
		if item.Layout != nil {
			// 布局UI
			element = append(element, item.Layout.Render())
		}
	}
	return element
}

// Render 渲染表单
func (t *Form) Render(ctx echo.Context) *node.TNode {

	if t.key > 0 && t.model != nil {
		// 预加载链表
		for _, item := range t.element {
			if item.HasAs != "" {
				t.modelDB.Preload(item.HasAs)
			}
		}
		// 查询当前数据
		queryModel := t.model
		t.modelDB.First(queryModel, t.key)
		jsonData, _ := json.Marshal(queryModel)
		_ = json.Unmarshal(jsonData, &t.info)
	}
	// 获取默认数据
	elements := t.ExpandElement()

	data := map[string]any{}
	for _, item := range elements {
		data[item.Field] = item.GetData(t.info)
	}
	// 渲染表单元素
	element := t.RenderElement()
	// 渲染表单UI
	formNode := &Node{
		Url:    t.url,
		Method: t.method,
		Title:  t.title,
	}
	formNode.SetHeader(t.header)
	formNode.SetElement(element)
	formNode.SetData(&data)
	formNode.SetBack(t.back)
	formNode.SetDialog(t.dialog)
	return formNode.Render()
}

// ExpandElement 展开元素
func (t *Form) ExpandElement() []*Element {
	var data []*Element
	for _, item := range t.element {
		data = append(data, item)
		if item.Layout != nil {
			data = append(data, item.Layout.Expand()...)
		}
	}
	return data
}

// SaveFn 自定义保存函数
func (t *Form) SaveFn(callback func(data map[string]any, key uint) error) {
	t.saveFn = callback
}

// SaveBefore 保存前处理
func (t *Form) SaveBefore(callback func(data map[string]any, update bool, db *gorm.DB) error) {
	t.saveBefore = callback
}

// SaveAfter 保存后处理
func (t *Form) SaveAfter(callback func(data any, update bool) error) {
	t.saveAfter = callback
}

// Save 保存表单
func (t *Form) Save(ctx echo.Context) error {
	var err error

	// 获取post字段
	postData := make(map[string]any)
	buf := function.CtxBody(ctx)
	err = json.Unmarshal(buf, &postData)
	if err != nil {
		panic("Unable to parse json data")
	}

	// 获取元素数据
	elements := t.ExpandElement()

	// 过滤提交字段
	data := map[string]any{}

	for _, item := range elements {
		if item.HasAs != "" {
			continue
		}
		if len(item.Verify) > 0 {
			for _, val := range item.Verify {
				err := t.validate.Var(postData[item.Field], val["role"])
				if err != nil {
					if val["message"] != "" {
						return exception2.ParameterError(val["message"])
					}
					return err
				}
			}
		}
		data[item.Field] = postData[item.Field]
	}

	// 格式化字符表单元素
	for _, item := range elements {
		value := data[item.Field]
		// 通过数据格式化
		if item.Format != nil {
			value = item.Format(value)
		}
		// 通过元素格式化保存
		value = item.SaveData(value, data)
		data[item.Field] = value
	}

	// 自定义保存
	if t.saveFn != nil {
		return t.saveFn(data, t.key)
	}

	// 非模型保存
	if t.model == nil {
		return nil
	}

	// 过滤字段
	fields := []string{}
	result, _ := global.Db.Migrator().ColumnTypes(t.model)
	for _, col := range result {
		fields = append(fields, col.Name())
	}
	for k, _ := range data {
		_, _, ok := lo.FindIndexOf[string](fields, func(i string) bool {
			return i == k
		})
		if !ok {
			delete(data, k)
		}
	}

	// 更新状态
	updateStatus := lo.Ternary[bool](t.key > 0, true, false)

	// 事务开启
	transaction := global.Db.Begin()

	// 保存前数据
	if t.saveBefore != nil {
		err = t.saveBefore(data, updateStatus, transaction)
		if err != nil {
			return err
		}
	}
	fmt.Println(data)

	// 获取树形字段
	ret := reflect.TypeOf(t.model).Elem()
	_, parentField := ret.FieldByName("ParentID")
	_, sortField := ret.FieldByName("Sort")
	orgModel := reflect.New(ret).Interface()
	if parentField && sortField {
		tmpId := cast.ToUint(data["parent_id"])
		if tmpId > 0 {
			var pNum int64
			transaction.Model(t.model).Where(t.primary+" = ?", tmpId).Count(&pNum)
			if pNum == 0 {
				return exception2.BusinessError("parent data does not exist")
			}
		}

	}

	// 保存基本数据
	mode := transaction.Model(t.model)
	if t.key > 0 {
		err = mode.Where(t.primary+" = ?", t.key).Updates(data).Error
	} else {
		err = mode.Create(data).Error
	}
	if err != nil {
		return exception2.Error(err)
	}
	if t.key == 0 {
		lastData := map[string]any{}
		err = transaction.Model(t.model).Select(t.primary).Last(&lastData).Error
		if err != nil {
			return exception2.Error(err)
		}
		t.key = cast.ToUint(lastData[t.primary])
	}

	// 查询数据
	err = transaction.Model(t.model).Find(t.model, t.key).Error
	if err != nil {
		return exception2.Error(err)
	}
	marshal, _ := json.Marshal(t.model)
	json.Unmarshal(marshal, &t.info)

	//排序处理
	if parentField && sortField {
		// 获取最大排序值
		pid := cast.ToUint(t.info["parent_id"])
		listData := map[string]any{}
		res := transaction.Model(&orgModel).Select("MAX(sort) as latest")
		if pid == 0 {
			res.Where("parent_id is NULL")
		} else {
			res.Where("parent_id = ?", pid)
		}
		res.Debug().Scan(&listData)
		if res.Error != nil {
			return err
		}
		// 更新顺序
		latest := cast.ToUint(listData["latest"]) + 1
		err = transaction.Model(&orgModel).Where(t.primary+" = ?", t.key).Update("sort", latest).Error
		if err != nil {
			return err
		}
	}

	// 多对多关联数据
	for _, item := range elements {
		if item.HasAs == "" {
			continue
		}
		model := transaction.Model(t.model).Association(item.HasAs)

		// 解析关联ID
		var hasIds []any
		switch postData[item.Field].(type) {
		case []any:
			hasIds = cast.ToSlice(postData[item.Field])
			break
		default:
			hasIds = []any{postData[item.Field]}
		}

		// 构建结构体
		var hasData []map[string]any
		for _, id := range hasIds {
			hasData = append(hasData, map[string]any{
				item.HasKey: cast.ToInt(id),
			})
		}
		err = mapstructure.Decode(hasData, item.HasModel)
		if err != nil {
			return err
		}
		// 替换当前关联
		err = model.Replace(item.HasModel)
		if err != nil {
			return err
		}
	}

	if t.saveAfter != nil {
		err = t.saveAfter(t.model, updateStatus)
		if err != nil {
			return err
		}
	}

	// 事务提交
	transaction.Commit()

	return err
}
