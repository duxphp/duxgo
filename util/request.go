package util

import (
	"github.com/duxphp/duxgo/v2/registry"
	"github.com/duxphp/duxgo/v2/validator"
	"github.com/labstack/echo/v4"
)

// RequestParser 请求解析验证
func RequestParser(ctx echo.Context, params any) error {
	var err error
	if err = ctx.Bind(params); err != nil {
		return err
	}
	err = registry.Validator.Struct(params)
	if err = validator.ProcessError(params, err); err != nil {
		return err
	}
	return nil
}

// QueryParser 请求解析验证
func QueryParser(ctx echo.Context, params any) error {
	var err error
	if err = ctx.Bind(params); err != nil {
		return err
	}
	err = registry.Validator.Struct(params)
	if err = validator.ProcessError(params, err); err != nil {
		return err
	}
	return nil
}
