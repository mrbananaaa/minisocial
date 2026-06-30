package httpx

import (
	"encoding/json"
	"net/http"
)

type Response[T, M any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
	Meta    M      `json:"meta,omitempty"`
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	payload any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}

// Res ...
// use this to enforce the Response struct so we have structured api response
// TODO: refactor this later
func Res[T any](
	w http.ResponseWriter,
	status int,
	message string,
	data T,
) {
	writeJSON(w, status, Response[T, any]{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ResWithMeta[T, M any](
	w http.ResponseWriter,
	status int,
	message string,
	data T,
	meta M,
) {
	writeJSON(w, status, Response[T, M]{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func OK[T any](w http.ResponseWriter, data T) {
	writeJSON(w, http.StatusOK, Response[T, any]{
		Success: true,
		Data:    data,
	})
}

func OKWithMeta[T, M any](w http.ResponseWriter, data T, meta M) {
	writeJSON(w, http.StatusOK, Response[T, M]{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

func Created[T any](w http.ResponseWriter, data T) {
	writeJSON(w, http.StatusCreated, Response[T, any]{
		Success: true,
		Data:    data,
	})
}

func Message(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, Response[any, any]{
		Success: status < 400,
		Message: msg,
	})
}
