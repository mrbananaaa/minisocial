package application

import "github.com/mrbananaaa/minisocial/internal/post/domain"

type Application struct {
	repo domain.Repository
}

func New(
	repo domain.Repository,
) *Application {
	return &Application{
		repo: repo,
	}
}
