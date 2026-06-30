package application

import "github.com/mrbananaaa/minisocial/internal/user/domain"

type Service struct {
	repo   domain.Repository
	hasher PasswordHasher
}

func New(
	repo domain.Repository,
	hasher PasswordHasher,
) *Service {
	return &Service{
		repo:   repo,
		hasher: hasher,
	}
}
