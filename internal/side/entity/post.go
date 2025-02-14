package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type Post struct {
	ID        ulid.ULID
	Body      string
	CreatedAt int64
	UpdatedAt int64
	Author    *Author
	Side      *Side
}

type Author struct {
	ID       uuid.UUID
	Username string
}

func CreatePost(body string, authorID uuid.UUID, sideID uuid.UUID) *Post {
	return &Post{
		ID:        ulid.Make(),
		Body:      body,
		CreatedAt: time.Now().UnixMilli(),
		Author:    &Author{ID: authorID},
		Side:      &Side{ID: sideID},
	}
}
