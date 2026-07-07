package createuser

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/platform/db"
	userApp "github.com/mrbananaaa/minisocial/internal/user/application"
	userDomain "github.com/mrbananaaa/minisocial/internal/user/domain"
)

type Workflow struct {
	users     UserService
	txManager *db.TxManager
}

type UserService interface {
	CreateUser(ctx context.Context, input userApp.CreateUserInput) (*userDomain.User, error)
}

func New(
	users UserService,
	txManager *db.TxManager,
) *Workflow {
	return &Workflow{
		users:     users,
		txManager: txManager,
	}
}

type Input struct {
	Email    string
	Username string
	Name     string
	Password string
}

type Output struct {
	ID        uuid.UUID
	Email     string
	Username  string
	Name      string
	Bio       string
	AvatarURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (w *Workflow) Execute(
	ctx context.Context,
	input Input,
) (*Output, error) {
	var out *Output

	err := w.txManager.WithTx(ctx, func(ctx context.Context) error {
		u, err := w.users.CreateUser(ctx, userApp.CreateUserInput(input))
		if err != nil {
			return err
		}

		out = &Output{
			ID:        u.ID,
			Email:     u.Email,
			Username:  u.Username,
			Name:      u.Name,
			Bio:       u.Bio,
			AvatarURL: u.AvatarURL,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}
