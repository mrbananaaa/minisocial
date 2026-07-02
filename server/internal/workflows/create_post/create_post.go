package createpost

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/mrbananaaa/minisocial/internal/platform/db"
	userDomain "github.com/mrbananaaa/minisocial/internal/user/domain"

	postApp "github.com/mrbananaaa/minisocial/internal/post/application"
	postDomain "github.com/mrbananaaa/minisocial/internal/post/domain"
)

type Workflow struct {
	users     UserService
	posts     PostService
	txManager *db.TxManager
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
	txManager *db.TxManager,
) *Workflow {
	return &Workflow{
		users:     users,
		posts:     posts,
		txManager: txManager,
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

// Execute ...
// TODO: Test this manually on the handler LOL 😂
func (w *Workflow) Execute(
	ctx context.Context,
	input Input,
) (*Output, error) {
	var out *Output

	err := w.txManager.WithTx(ctx, func(ctx context.Context) error {
		u, err := w.users.GetUserByID(ctx, input.AuthorID)
		if err != nil {
			return err
		}

		p, err := w.posts.CreatePost(ctx, postApp.CreatePostInput{
			AuthorID: u.ID,
			Title:    input.Title,
			Content:  input.Content,
		})
		if err != nil {
			return err
		}

		out = &Output{
			ID:        p.ID,
			AuthorID:  p.AuthorID,
			Title:     p.Title,
			Slug:      p.Slug,
			Content:   p.Content,
			Status:    p.Status,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}
