package httpx

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mrbananaaa/minisocial/internal/platform/validation"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func Error(
	w http.ResponseWriter,
	status int,
	message string,
	details any,
) {
	res := &ErrorResponse{
		Success: false,
		Message: message,
	}

	if details != nil {
		res.Details = details
	}

	writeJSON(w, status, res)
}

func ErrInvalidUUID(w http.ResponseWriter, field string) {
	Error(
		w,
		http.StatusBadRequest,
		fmt.Sprintf("invalid uuid field - %s", field),
		nil,
	)
}

func ErrInvalidRequestBody(w http.ResponseWriter) {
	writeJSON(w, http.StatusBadRequest, ErrorResponse{
		Success: false,
		Message: "invalid request body",
	})
}

func ErrValidation(w http.ResponseWriter, err error) {
	if err == nil {
		Error(
			w,
			http.StatusInternalServerError,
			"something went wrong",
			nil,
		)
		return
	}

	var validationError *validation.ValidationError
	if errors.As(err, &validationError) {
		Error(
			w,
			http.StatusBadRequest,
			validationError.Error(),
			validationError.Errors,
		)
		return
	}

	Error(
		w,
		http.StatusInternalServerError,
		"something went wrong",
		nil,
	)
}
