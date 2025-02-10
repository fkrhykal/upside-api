package entity

import (
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
