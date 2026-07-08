// Package api post transport layer
package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/platform/httpx"
	"github.com/mrbananaaa/minisocial/internal/platform/validation"
	"github.com/mrbananaaa/minisocial/internal/post/application"
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
		app:        app,
		createpost: createpost,
		validator:  validator,
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

	// db60cc82-3ff6-4c89-a16b-60c340f18f3c admin
	// bc481a67-6272-4137-a449-4890d6d76b34 another admin
	// 852b1ddb-2670-47b4-a0d7-7480c5861687 yet another admin
	userID, _ := uuid.Parse("852b1ddb-2670-47b4-a0d7-7480c5861687")
	out, err := h.createpost.Execute(r.Context(), createpost.Input{
		// TODO: extract user from context
		AuthorID: userID,
		Title:    req.Title,
		Content:  req.Content,
	})
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	httpx.Res(w, http.StatusCreated, "post created", PostResponse{
		ID:        out.ID,
		AuthorID:  out.AuthorID,
		Title:     out.Title,
		Slug:      out.Slug,
		Content:   out.Content,
		Status:    string(out.Status),
		CreatedAt: out.CreatedAt,
		UpdatedAt: out.UpdatedAt,
	})
}

func (h *Handler) EditPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		httpx.ErrInvalidUUID(w, "post_id param")
		return
	}

	var req EditPostRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrInvalidRequestBody(w)
		return
	}

	if err = h.validator.Validate(req); err != nil {
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
