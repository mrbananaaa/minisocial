package domain

import "errors"

var (
	ErrPostArchived          = errors.New("post is archived")
	ErrPostTitleLength       = errors.New("title length must be at least 6 chars long and max 255 chars long")
	ErrPostContentLength     = errors.New("content length must be at least 20 chars long")
	ErrPostNotFound          = errors.New("post not found")
	ErrPostSlugAlreadyExists = errors.New("post slugs already exists")
)
