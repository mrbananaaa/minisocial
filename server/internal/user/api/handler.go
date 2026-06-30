package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/platform/httpx"
	"github.com/mrbananaaa/minisocial/internal/platform/validation"
	"github.com/mrbananaaa/minisocial/internal/user/application"
)

type Handler struct {
	app       *application.Application
	validator *validation.Validator
}

func New(
	app *application.Application,
	validator *validation.Validator,
) *Handler {
	return &Handler{
		app:       app,
		validator: validator,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrInvalidRequestBody(w)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpx.ErrValidation(w, err)
		return
	}

	out, err := h.app.CreateUser(r.Context(), toCreateUserInput(req))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	httpx.Res(w, http.StatusCreated, "user created", toUserResponse(out.User))
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		httpx.ErrInvalidUUID(w, "user_id param")
		return
	}

	user, err := h.app.GetUser(r.Context(), application.GetUserInput{
		ID: new(userID),
	})
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	httpx.Res(w, http.StatusOK, "", toUserResponse(user))
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		httpx.ErrInvalidUUID(w, "user_id param")
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrInvalidRequestBody(w)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpx.ErrValidation(w, err)
		return
	}

	updated, err := h.app.UpdateProfile(r.Context(), toUpdateUserInput(userID, req))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	httpx.Res(w, http.StatusOK, "user updated", toUserResponse(updated))
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		httpx.ErrInvalidUUID(w, "user_id param")
		return
	}

	err = h.app.DeleteUser(r.Context(), userID)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	httpx.Res[any](w, http.StatusOK, "user deleted", nil)
}
