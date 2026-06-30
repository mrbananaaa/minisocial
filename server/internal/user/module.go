package user

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/mrbananaaa/minisocial/internal/platform/validation"
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
	validator := validation.New()
	app := application.New(repo, hasher)
	handler := api.New(app, validator)

	return &Module{
		app:     app,
		handler: handler,
	}
}

func (m *Module) Service() *application.Application {
	return m.app
}

func (m *Module) RegisterRoutes(r chi.Router) {
	m.handler.RegisterRoutes(r)
}
