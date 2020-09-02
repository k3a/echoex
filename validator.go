package echoex

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator"
)

// CustomValidator implements echo validator using go-playground's validator.
// It allows to use validate tags like `validate:"required,email"`.
// For the list of all supported tags see https://godoc.org/gopkg.in/go-playground/validator.v9
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator creates a new CustomValidator instance just encapsulating an existing
// validator.Validate without any extra functionality
func NewCustomValidatorUsingValidate(validate *validator.Validate) *CustomValidator {
	return &CustomValidator{validator: validate}
}

// NewCustomValidator creates a new CustomValidator instance encapsulating
// a new validator.Validate instance configured to
func NewCustomValidator() *CustomValidator {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name != "-" && name != "" {
			return name
		}
		name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		if name != "-" && name != "" {
			return name
		}
		name = strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]
		if name != "-" && name != "" {
			return name
		}
		return ""
	})

	return NewCustomValidatorUsingValidate(v)
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
