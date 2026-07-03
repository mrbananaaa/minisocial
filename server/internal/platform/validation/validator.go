// Package validation contain validation implementation aka go-playground/validator wrapper
package validation

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New(opts ...Option) *Validator {
	v := &Validator{
		validate: validator.New(
			validator.WithRequiredStructEnabled(),
		),
	}

	for _, opt := range opts {
		opt(v)
	}

	v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name, _, _ := strings.Cut(fld.Tag.Get("json"), ",")

		if name == "-" {
			return ""
		}

		return name
	})

	return v
}

func (v *Validator) Validate(val any) error {
	err := v.validate.Struct(val)
	if err == nil {
		return nil
	}

	if validationErrors, ok := errors.AsType[validator.ValidationErrors](err); ok {
		return NewValidationError(validationErrors)
	}

	return err
}
