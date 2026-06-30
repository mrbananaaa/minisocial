package api

import (
	"errors"
	"net/http"

	"github.com/mrbananaaa/minisocial/internal/user/domain"
)

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		http.Error(w, err.Error(), http.StatusConflict)

	case errors.Is(err, domain.ErrUsernameAlreadyExists):
		http.Error(w, err.Error(), http.StatusConflict)

	case errors.Is(err, domain.ErrUserNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
