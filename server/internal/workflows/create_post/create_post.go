package createpost

import (
	"context"
	"time"

	"github.com/google/uuid"

	userDomain "github.com/mrbananaaa/minisocial/internal/user/domain"

	postApp "github.com/mrbananaaa/minisocial/internal/post/application"
	postDomain "github.com/mrbananaaa/minisocial/internal/post/domain"
)

type Workflow struct {
	users UserService
	posts PostService
}

type UserService interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*userDomain.User, error)
}

type PostService interface {
	CreatePost(ctx context.Context, input postApp.CreatePostInput) (*postDomain.Post, error)
}

func New(
	users UserService,
	posts PostService,
) *Workflow {
	return &Workflow{
		users: users,
		posts: posts,
	}
}

type Input struct {
	AuthorID uuid.UUID
	Title    string
	Content  string
}

type Output struct {
	ID        uuid.UUID
	AuthorID  uuid.UUID
	Title     string
	Slug      string
	Content   string
	Status    postDomain.Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (w *Workflow) Execute(
	ctx context.Context,
	input Input,
) (*Output, error) {
	u, err := w.users.GetUserByID(ctx, input.AuthorID)
	if err != nil {
		return nil, err
	}

	p, err := w.posts.CreatePost(ctx, postApp.CreatePostInput{
		AuthorID: u.ID,
		Title:    input.Title,
		Content:  input.Content,
	})
	if err != nil {
		return nil, err
	}

	return &Output{
		ID:        p.ID,
		AuthorID:  p.AuthorID,
		Title:     p.Title,
		Slug:      p.Slug,
		Content:   p.Content,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}, nil
}
