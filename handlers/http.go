package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// New Http 异常
func New(code int, msg ...string) *fiber.Error {
	return fiber.NewError(code, msg...)
}

// NotFound 页面不存在
func NotFound() *fiber.Error {
	return fiber.NewError(fiber.StatusNotFound)
}

// BusinessError 业务错误
func BusinessError(msg ...string) *fiber.Error {
	return fiber.NewError(
		http.StatusInternalServerError,
		msg...,
	)
}

// BusinessErrorf 业务错误
func BusinessErrorf(msg string, params ...any) *fiber.Error {
	return fiber.NewError(
		http.StatusInternalServerError,
		fmt.Sprintf(msg, params),
	)
}

// ParameterError 参数错误
func ParameterError(msg ...string) *fiber.Error {
	return fiber.NewError(
		http.StatusBadRequest,
		msg...,
	)
}

// ParameterErrorf 参数错误
func ParameterErrorf(msg string, params ...any) *fiber.Error {
	return fiber.NewError(
		http.StatusBadRequest,
		fmt.Sprintf(msg, params),
	)
}

// UnknownError 未知错误
func UnknownError(msg ...string) *fiber.Error {
	return fiber.NewError(
		http.StatusForbidden,
		msg...,
	)
}

// UnknownErrorf 未知错误
func UnknownErrorf(msg string, params ...any) *fiber.Error {
	return fiber.NewError(
		http.StatusForbidden,
		fmt.Sprintf(msg, params),
	)
}
