package application

import "github.com/mrbananaaa/minisocial/internal/user/domain"

type Application struct {
	repo   domain.Repository
	hasher PasswordHasher
}

func New(
	repo domain.Repository,
	hasher PasswordHasher,
) *Application {
	return &Application{
		repo:   repo,
		hasher: hasher,
	}
}
