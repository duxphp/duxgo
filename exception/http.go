package exception

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func New(code int, msg ...any) *echo.HTTPError {
	return echo.NewHTTPError(code, msg...)
}

// NotFound 页面不存在
func NotFound() *echo.HTTPError {
	return New(
		http.StatusNotFound,
		http.StatusText(http.StatusNotFound),
	)
}

// BusinessError 业务错误
func BusinessError(msg string) *echo.HTTPError {
	return New(
		http.StatusInternalServerError,
		msg,
	)
}

// BusinessErrorf 业务错误
func BusinessErrorf(msg string, params ...any) *echo.HTTPError {
	return New(
		http.StatusInternalServerError,
		fmt.Sprintf(msg, params),
	)
}

// ParameterError 参数错误
func ParameterError(msg string) *echo.HTTPError {
	return New(
		http.StatusBadRequest,
		msg,
	)
}

// ParameterErrorf 参数错误
func ParameterErrorf(msg string, params ...any) *echo.HTTPError {
	return New(
		http.StatusBadRequest,
		fmt.Sprintf(msg, params),
	)
}

// UnknownError 未知错误
func UnknownError(msg string) *echo.HTTPError {
	return New(
		http.StatusForbidden,
		msg,
	)
}

// UnknownErrorf 未知错误
func UnknownErrorf(msg string, params ...any) *echo.HTTPError {
	return New(
		http.StatusForbidden,
		fmt.Sprintf(msg, params),
	)
}
