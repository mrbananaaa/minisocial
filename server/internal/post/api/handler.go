package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/platform/httpx"
	"github.com/mrbananaaa/minisocial/internal/platform/validation"
	"github.com/mrbananaaa/minisocial/internal/post/application"
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

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req CreatePostRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrInvalidRequestBody(w)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpx.ErrValidation(w, err)
		return
	}

	out, err := h.app.CreatePost(r.Context(), toCreatePostInput(uuid.New(), req))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	httpx.Res(w, http.StatusCreated, "post created", toPostResponse(out))
}

func (h *Handler) EditPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		httpx.ErrInvalidUUID(w, "post_id param")
		return
	}

	var req EditPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrInvalidRequestBody(w)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		httpx.ErrValidation(w, err)
		return
	}

	out, err := h.app.EditPost(r.Context(), toEditPostInput(postID, req))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	httpx.Res(w, http.StatusOK, "post updated", toPostResponse(out))
}

func (h *Handler) ArchivePost(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		httpx.ErrInvalidUUID(w, "post_id param")
		return
	}

	if err := h.app.ArchivePost(r.Context(), postID); err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	httpx.Message(w, http.StatusOK, "post archived")
}
