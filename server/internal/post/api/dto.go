package api

import (
	"time"

	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/post/application"
	"github.com/mrbananaaa/minisocial/internal/post/domain"
)

type PostResponse struct {
	ID         uuid.UUID  `json:"id"`
	AuthorID   uuid.UUID  `json:"author_id"`
	Title      string     `json:"title"`
	Slug       string     `json:"slug"`
	Content    string     `json:"content"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at,omitempty"`
}

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,min=6,max=255"`
	Content string `json:"content" validate:"required,min=6"`
}

type EditPostRequest struct {
	Title   *string `json:"title,omitempty" validate:"min=6,max=255"`
	Content *string `json:"content,omitempty" validate:"min=6"`
}

func toCreatePostInput(authordID uuid.UUID, req CreatePostRequest) application.CreatePostInput {
	return application.CreatePostInput{
		AuthorID: authordID,
		Title:    req.Title,
		Content:  req.Content,
	}
}

func toEditPostInput(id uuid.UUID, req EditPostRequest) application.EditPostInput {
	return application.EditPostInput{
		ID:      id,
		Title:   req.Title,
		Content: req.Content,
	}
}

func toPostResponse(post *domain.Post) PostResponse {
	return PostResponse{
		ID:         post.ID,
		AuthorID:   post.AuthorID,
		Title:      post.Title,
		Slug:       post.Slug,
		Content:    post.Content,
		Status:     string(post.Status),
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
		ArchivedAt: post.ArchivedAt,
	}
}
