package request

import (
	"github.com/duxphp/duxgo/v2/validator"
	"github.com/gofiber/fiber/v2"
)

// Parser 请求解析验证
func Parser(ctx *fiber.Ctx, params any) error {
	var err error
	if err = ctx.BodyParser(params); err != nil {
		return err
	}
	err = validator.Validator().Struct(params)
	if err = validator.ProcessError(params, err); err != nil {
		return err
	}
	return nil
}
