package post

import (
	"github.com/mrbananaaa/minisocial/internal/post/application"
	"github.com/mrbananaaa/minisocial/internal/post/repositories/postgres"
	"github.com/mrbananaaa/minisocial/internal/post/repositories/postgres/sqlc"
)

type Module struct {
	app *application.Application
}

func New(
	db sqlc.DBTX,
) *Module {
	repo := postgres.New(db)
	app := application.New(repo)

	return &Module{
		app: app,
	}
}

func (m *Module) Service() *application.Application {
	return m.app
}
