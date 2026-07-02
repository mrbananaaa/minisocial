package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/platform/httpx"
	"github.com/mrbananaaa/minisocial/internal/platform/validation"
	"github.com/mrbananaaa/minisocial/internal/post/application"
	"github.com/mrbananaaa/minisocial/internal/post/domain"
	createpost "github.com/mrbananaaa/minisocial/internal/workflows/create_post"
)

type Handler struct {
	app        *application.Application
	createpost *createpost.Workflow
	validator  *validation.Validator
}

func New(
	app *application.Application,
	createpost *createpost.Workflow,
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

	out, err := h.createpost.Execute(r.Context(), createpost.Input{
		// TODO: extract user from context
		AuthorID: uuid.New(),
		Title:    req.Title,
		Content:  req.Content,
	})
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	httpx.Res(w, http.StatusCreated, "post created", toPostResponse(&domain.Post{
		ID:        out.ID,
		AuthorID:  out.AuthorID,
		Title:     out.Title,
		Slug:      out.Slug,
		Content:   out.Content,
		Status:    out.Status,
		CreatedAt: out.CreatedAt,
		UpdatedAt: out.UpdatedAt,
	}))
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
