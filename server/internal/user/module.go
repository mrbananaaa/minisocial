package user

import (
	"log/slog"

	"github.com/mrbananaaa/minisocial/internal/user/api"
	"github.com/mrbananaaa/minisocial/internal/user/application"
	"github.com/mrbananaaa/minisocial/internal/user/infra/password"
	"github.com/mrbananaaa/minisocial/internal/user/repositories/postgres"
	"github.com/mrbananaaa/minisocial/internal/user/repositories/postgres/sqlc"
)

type Module struct {
	app     *application.Application
	handler *api.Handler
}

func New(
	db sqlc.DBTX,
	logger *slog.Logger,
) *Module {
	repo := postgres.New(db)
	hasher := password.New(12)
	app := application.New(repo, hasher)

	return &Module{
		app: app,
	}
}

func (m *Module) Service() *application.Application {
	return m.app
}
