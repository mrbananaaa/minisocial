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
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

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

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return NewValidationError(validationErrors)
	}

	return err
}
