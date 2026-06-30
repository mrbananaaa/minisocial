package workflows

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/platform/httpx"
	"github.com/mrbananaaa/minisocial/internal/platform/validation"
	createpost "github.com/mrbananaaa/minisocial/internal/workflows/create_post"

	postApp "github.com/mrbananaaa/minisocial/internal/post/application"
	userApp "github.com/mrbananaaa/minisocial/internal/user/application"
)

type Module struct {
	createPost *createpost.Workflow

	validator *validation.Validator
}

func New(
	users *userApp.Application,
	posts *postApp.Application,
) *Module {
	cp := createpost.New(users, posts)

	validator := validation.New()

	return &Module{
		createPost: cp,
		validator:  validator,
	}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Post("/posts", m.createPostHandler)
}

type CreatePostResponse struct {
	ID        uuid.UUID `json:"id"`
	AuthorID  uuid.UUID `json:"author_id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePostReq struct {
	Title   string `json:"title" validate:"required,min=6,max=255"`
	Content string `json:"content" validate:"required,min=6"`
}

func (m *Module) createPostHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: get user from context
	userID := uuid.New()

	var req CreatePostReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.ErrInvalidRequestBody(w)
		return
	}

	if err := m.validator.Validate(req); err != nil {
		httpx.ErrValidation(w, err)
		return
	}

	out, err := m.createPost.Execute(r.Context(), createpost.Input{
		AuthorID: userID,
		Title:    req.Title,
		Content:  req.Content,
	})
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	httpx.Res(w, http.StatusCreated, "post created", CreatePostResponse{
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
