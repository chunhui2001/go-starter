package commons

import (
	"errors"
	_ "github.com/chunhui2001/go-starter/utils"
	"github.com/go-playground/validator/v10"
)

func GetErrorMsg(fe validator.FieldError) string {

	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	}

	return "Unknown error"

}

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"msg"`
}

type R struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   error       `json:"-"`
}

func (r R) Get() map[string]interface{} {
	return Result(r)
}

func (r R) Success() map[string]interface{} {
	r.Code = 200
	return Result(r)
}

func (r R) Fail(code int) map[string]interface{} {
	r.Code = code
	return Result(r)
}

func Result(r R) map[string]interface{} {

	m := make(map[string]interface{})

	if r.Code != 0 {
		m["code"] = r.Code
	}

	if r.Message != "" {
		m["message"] = r.Message
	} else {
		m["message"] = "Ok."
	}

	if r.Data != nil {
		m["data"] = r.Data
	}

	if r.Error != nil {

		var ve validator.ValidationErrors

		if errors.As(r.Error, &ve) {
			out := make([]ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMsg{fe.Field(), GetErrorMsg(fe)}
			}
			m["message"] = "Validator-Failed."
			m["errors"] = out
		} else {
			m["message"] = r.Error.Error()
		}

	}

	return m
}
