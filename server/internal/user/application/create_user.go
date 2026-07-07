package application

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/user/domain"
)

type CreateUserInput struct {
	Email    string
	Username string
	Name     string
	Password string
}

func (a *Application) CreateUser(ctx context.Context, input CreateUserInput) (*domain.User, error) {
	exists, err := a.repo.GetByEmail(ctx, input.Email)
	if err != nil && errors.Is(err, domain.ErrEmailAlreadyExists) {
		return nil, err
	}
	if exists != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	exists, err = a.repo.GetByUsername(ctx, input.Username)
	if err != nil && errors.Is(err, domain.ErrUsernameAlreadyExists) {
		return nil, err
	}
	if exists != nil {
		return nil, domain.ErrUsernameAlreadyExists
	}

	if err = validatePassword(input.Password); err != nil {
		return nil, err
	}

	hashed, err := a.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        input.Email,
		Username:     input.Username,
		Name:         input.Name,
		PasswordHash: hashed,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := a.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func validatePassword(password string) error {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) < 6 {
		hasMinLen = true
	}

	if hasMinLen {
		return errors.New("password must be at least 6 characters long")
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one capital letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one non-capital letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
