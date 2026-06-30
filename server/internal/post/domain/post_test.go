package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mrbananaaa/minisocial/internal/post/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPost(t *testing.T) {
	const (
		title   = "My First Post"
		content = "Hello, MiniSocial!"
	)

	authorID := uuid.New()

	t.Run("success - new post w all expected value", func(t *testing.T) {
		p, err := domain.New(domain.NewPostInput{
			AuthorID: authorID,
			Title:    title,
			Content:  content,
		})
		require.NoError(t, err)

		assert.NotEqual(t, uuid.Nil, p.ID)
		assert.Equal(t, authorID, p.AuthorID)
		assert.Equal(t, title, p.Title)
		assert.Equal(t, "", p.Slug)
		assert.Equal(t, content, p.Content)
		assert.Equal(t, domain.StatusDraft, p.Status)
		assert.False(t, p.CreatedAt.IsZero())
		assert.False(t, p.UpdatedAt.IsZero())
		assert.Equal(t, p.CreatedAt, p.UpdatedAt)
		assert.Nil(t, p.ArchivedAt)
	})

	t.Run("fail - empty title", func(t *testing.T) {
		p, err := domain.New(domain.NewPostInput{
			AuthorID: authorID,
			Title:    "",
			Content:  content,
		})

		require.Nil(t, p)
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrPostTitleLength)
	})

	t.Run("fail - empty content", func(t *testing.T) {
		authorID := uuid.New()

		p, err := domain.New(domain.NewPostInput{
			AuthorID: authorID,
			Title:    title,
			Content:  "",
		})

		require.Nil(t, p)
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrPostContentLength)
	})

	t.Run("fail - can't update archived post", func(t *testing.T) {
		p, err := domain.New(domain.NewPostInput{
			AuthorID: authorID,
			Title:    title,
			Content:  content,
		})
		require.NoError(t, err)

		err = p.Archive()
		require.NoError(t, err)

		err = p.Edit("edited title", "edited content")
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrPostArchived)
	})
}
