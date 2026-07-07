// Package domain post
package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/platform/events"
)

type Status string

const (
	StatusDraft Status = "draft"
	// StatusPublished Status = "published"
	StatusArchived Status = "archived"
)

type Post struct {
	ID         uuid.UUID
	AuthorID   uuid.UUID
	Title      string
	Slug       string
	Content    string
	Status     Status
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ArchivedAt *time.Time

	events []events.Event
}

func (p *Post) addEvent(event events.Event) {
	p.events = append(p.events, event)
}

func (p *Post) PullEvents() []events.Event {
	events := p.events
	p.events = nil

	return events
}

type NewPostInput struct {
	AuthorID uuid.UUID
	Title    string
	Content  string
}

func New(input NewPostInput) (*Post, error) {
	if len(input.Title) <= 6 || len(input.Title) > 255 {
		return nil, ErrPostTitleLength
	}

	if len(input.Content) <= 6 {
		return nil, ErrPostContentLength
	}

	now := time.Now()
	post := &Post{
		ID:         uuid.New(),
		AuthorID:   input.AuthorID,
		Title:      input.Title,
		Slug:       uuid.NewString(), // generate slug later
		Content:    input.Content,
		Status:     StatusDraft,
		CreatedAt:  now,
		UpdatedAt:  now,
		ArchivedAt: nil,
	}

	post.addEvent(PostCreated{
		PostID:   post.ID,
		AuthorID: post.AuthorID,
		Title:    post.Title,
	})

	return post, nil
}

func (p *Post) Edit(title, content string) error {
	if p.Status == StatusArchived {
		return ErrPostArchived
	}

	if strings.TrimSpace(title) != "" {
		p.Title = title
	}

	if strings.TrimSpace(content) != "" {
		p.Content = content
	}

	p.UpdatedAt = time.Now()

	return nil
}

func (p *Post) Archive() error {
	if p.Status == StatusArchived {
		return ErrPostArchived
	}

	now := time.Now()

	p.Status = StatusArchived
	p.ArchivedAt = &now
	p.UpdatedAt = now

	return nil
}

// func (p *Post) Publish() {}
