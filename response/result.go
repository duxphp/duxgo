package response

import (
	"github.com/gofiber/fiber/v2"
)

// Result 返回数据
type Result struct {
	Ctx *fiber.Ctx
}

// New 构建消息
func New(ctx *fiber.Ctx) *Result {
	return &Result{Ctx: ctx}
}

// ResultData 数据结构体
type ResultData struct {
	Code    int    `json:"code" example:"200"`   //提示代码
	Message string `json:"message" example:"ok"` //提示信息
	Data    any    `json:"data"`                 //数据
}

// Render 模板渲染
func (r *Result) Render(name string, bind any) error {
	return r.Ctx.Render(name, bind)
}

// Send 发送消息
func (r *Result) Send(message string, data ...any) error {
	var params any
	if len(data) > 0 {
		params = data[0]
	} else {
		params = map[string]any{}
	}
	res := ResultData{}
	res.Code = 200
	res.Message = message
	res.Data = params
	return r.Ctx.Status(fiber.StatusOK).JSON(res)
}
