package dto

import (
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type CreatePostRequest struct {
	SideID uuid.UUID `json:"-"`
	Body   string    `json:"body" validate:"required"`
}

type CreatePostResponse struct {
	ID ulid.ULID `json:"id"`
}
