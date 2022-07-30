package validator

import (
	"errors"
	"github.com/duxphp/duxgo/core"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
)

func Init() {
	core.Validator = validator.New()
	err := core.Validator.RegisterValidation("cnPhone", func(f validator.FieldLevel) bool {
		value := f.Field().String()
		result, _ := regexp.MatchString(`^(1\d{10})$`, value)
		return result
	})
	if err != nil {
		return
	}
}

// ProcessError 验证错误处理
func ProcessError(object any, err error) error {
	if err == nil {
		return nil
	}
	invalid, ok := err.(*validator.InvalidValidationError)
	if ok {
		return errors.New("参数错误：" + invalid.Error())
	}
	validationErrs := err.(validator.ValidationErrors)
	for _, item := range validationErrs {
		fieldName := item.Field()
		typeOf := reflect.TypeOf(object)
		if typeOf.Kind() == reflect.Ptr {
			typeOf = typeOf.Elem()
		}
		field, ok := typeOf.FieldByName(fieldName)
		if ok {
			msg := field.Tag.Get("validateMsg")
			if msg != "" {
				return errors.New(msg)
			} else {
				return errors.New(item.Error())
			}

		} else {
			return errors.New(item.Error())
		}
	}
	return nil
}
