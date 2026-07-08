// Package createpost contain creating post workflow logic
package createpost

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/mrbananaaa/minisocial/internal/platform/db"
	"github.com/mrbananaaa/minisocial/internal/platform/events"
	"github.com/mrbananaaa/minisocial/internal/platform/outbox"
	userDomain "github.com/mrbananaaa/minisocial/internal/user/domain"

	postApp "github.com/mrbananaaa/minisocial/internal/post/application"
	postDomain "github.com/mrbananaaa/minisocial/internal/post/domain"
)

type Workflow struct {
	users      UserService
	posts      PostService
	txManager  *db.TxManager
	outboxRepo outbox.Repository
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
	outboxRepo outbox.Repository,
) *Workflow {
	return &Workflow{
		users:      users,
		posts:      posts,
		txManager:  txManager,
		outboxRepo: outboxRepo,
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
	var out *Output

	collector := events.NewCollector()
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
		collector.Track(p)

		events := collector.Flush()
		for _, event := range events {
			evt := outbox.NewEvent(
				event.AggregateType(),
				event.AggregateID(),
				event.EventType(),
				event.Payload(),
			)

			if err := w.outboxRepo.InsertEvent(ctx, evt); err != nil {
				return err
			}
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
