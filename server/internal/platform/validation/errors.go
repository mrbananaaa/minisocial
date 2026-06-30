package validation

import "github.com/go-playground/validator/v10"

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError struct {
	Errors []FieldError `json:"errors"`
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

func NewValidationError(errs validator.ValidationErrors) *ValidationError {
	var result []FieldError

	for _, e := range errs {
		result = append(result, FieldError{
			Field:   e.Field(),
			Message: messageForTag(e),
		})
	}

	return &ValidationError{Errors: result}
}

func messageForTag(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return "is too short"
	case "max":
		return "is too long"
	case "gte":
		return "must be greater or equal to " + e.Param()
	}

	return "is invalid"
}
