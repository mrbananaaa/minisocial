package post

import (
	"log/slog"

	"github.com/mrbananaaa/minisocial/internal/platform/validation"
	"github.com/mrbananaaa/minisocial/internal/post/api"
	"github.com/mrbananaaa/minisocial/internal/post/application"
	"github.com/mrbananaaa/minisocial/internal/post/repositories/postgres"
	"github.com/mrbananaaa/minisocial/internal/post/repositories/postgres/sqlc"
)

type Module struct {
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
		handler: handler,
	}
}
