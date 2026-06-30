package post

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/mrbananaaa/minisocial/internal/platform/validation"
	"github.com/mrbananaaa/minisocial/internal/post/api"
	"github.com/mrbananaaa/minisocial/internal/post/application"
	"github.com/mrbananaaa/minisocial/internal/post/repositories/postgres"
	"github.com/mrbananaaa/minisocial/internal/post/repositories/postgres/sqlc"
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
	validator := validation.New()
	app := application.New(repo)
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
