package exception

import (
	"fmt"
	"github.com/duxphp/duxgo/v2/registry"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type CoreError struct {
	Message string
}

func (e *CoreError) Error() string {
	return e.Message
}

// Error 外部错误
func Error(err any, params ...any) *CoreError {
	msg := "unknown error"
	if e, ok := err.(error); ok {
		msg = e.Error()
	} else if e, ok := err.(string); ok {
		msg = fmt.Sprintf(e, params)
	} else {
		msg = cast.ToString(err)
	}
	errs := &CoreError{
		Message: msg,
	}
	registry.Logger.Error().CallerSkipFrame(2).Interface("err", errs).Msg("core")
	return errs
}

// Internal 内部错误
func Internal(err any, params ...any) *CoreError {
	errs := Error(err, params)
	errs.Message = lo.Ternary[string](registry.DebugMsg == "", "business is busy, please try again", registry.DebugMsg)
	return errs
}
