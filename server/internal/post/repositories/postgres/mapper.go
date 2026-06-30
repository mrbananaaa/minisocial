package postgres

import (
	"github.com/mrbananaaa/minisocial/internal/post/domain"
	"github.com/mrbananaaa/minisocial/internal/post/repositories/postgres/sqlc"
)

func toDomain(p sqlc.Post) *domain.Post {
	return &domain.Post{
		ID:         p.ID,
		AuthorID:   p.AuthorID,
		Title:      p.Title,
		Slug:       p.Slug,
		Content:    p.Content,
		Status:     mapStatus(p.Status),
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		ArchivedAt: p.ArchivedAt,
	}
}

type FromDomainTypes interface {
	sqlc.CreatePostParams | sqlc.UpdatePostParams
}

// fromDomain ...
// It's totally bullshit and I know, I just wanna learn about generic.
func fromDomain[T FromDomainTypes](p *domain.Post) T {
	var zero T

	createParam := sqlc.CreatePostParams{
		ID:         p.ID,
		AuthorID:   p.AuthorID,
		Title:      p.Title,
		Slug:       p.Slug,
		Content:    p.Content,
		Status:     string(p.Status),
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		ArchivedAt: p.ArchivedAt,
	}
	updateParam := sqlc.UpdatePostParams{
		ID:         p.ID,
		Title:      p.Title,
		Slug:       p.Slug,
		Content:    p.Content,
		Status:     string(p.Status),
		UpdatedAt:  p.UpdatedAt,
		ArchivedAt: p.ArchivedAt,
	}

	switch any(zero).(type) {
	case sqlc.CreatePostParams:
		return any(createParam).(T)
	case sqlc.UpdatePostParams:
		return any(updateParam).(T)
	default:
		return any(createParam).(T)
	}
}

func mapStatus(s string) domain.Status {
	switch s {
	case "draft":
		return domain.StatusDraft
	case "archived":
		return domain.StatusArchived
	default:
		return domain.StatusDraft
	}
}
