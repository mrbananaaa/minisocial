// Package domain user
package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/platform/events"
)

type User struct {
	ID           uuid.UUID
	Email        string
	Username     string
	Name         string
	PasswordHash string
	Bio          string
	AvatarURL    string
	CreatedAt    time.Time
	UpdatedAt    time.Time

	events []events.Event
}

func (u *User) addEvent(event events.Event) {
	u.events = append(u.events, event)
}

func (u *User) PullEvents() []events.Event {
	events := u.events
	u.events = nil

	return events
}

type NewUserInput struct {
	Email        string
	Username     string
	Name         string
	PasswordHash string
}

func New(input NewUserInput) (*User, error) {
	now := time.Now()
	user := &User{
		ID:           uuid.New(),
		Email:        input.Email,
		Username:     input.Username,
		Name:         input.Name,
		PasswordHash: input.PasswordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	user.addEvent(UserCreated{
		UserID: user.ID,
		Email:  user.Email,
	})

	return user, nil
}

// TODO: validate input
func (u *User) UpdateName(name string) error {
	u.Name = name
	u.UpdatedAt = time.Now()

	return nil
}

// TODO: validate input
func (u *User) UpdateBio(bio string) error {
	u.Bio = bio
	u.UpdatedAt = time.Now()

	return nil
}

// TODO: validate input
func (u *User) UpdateAvatarURL(avatarURL string) error {
	u.AvatarURL = avatarURL
	u.UpdatedAt = time.Now()

	return nil
}
