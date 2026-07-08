// Package api user transport layer
package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/platform/httpx"
	"github.com/mrbananaaa/minisocial/internal/platform/validation"
	"github.com/mrbananaaa/minisocial/internal/user/application"
	createuser "github.com/mrbananaaa/minisocial/internal/workflows/create_user"
)

type Handler struct {
	app        *application.Application
	createuser *createuser.Workflow
	validator  *validation.Validator
}

func New(
	app *application.Application,
	createuser *createuser.Workflow,
	validator *validation.Validator,
) *Handler {
	return &Handler{
		app:        app,
		validator:  validator,
		createuser: createuser,
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

	user, err := h.createuser.Execute(r.Context(), createuser.Input{
		Email:    req.Email,
		Username: req.Username,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	httpx.Res(w, http.StatusCreated, "user created", UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Name:      user.Name,
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
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
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrInvalidRequestBody(w)
		return
	}

	if err = h.validator.Validate(req); err != nil {
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
