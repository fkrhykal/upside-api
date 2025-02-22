package dto

import (
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type VoteRequest struct {
	PostID   ulid.ULID
	VoteKind entity.VoteKind
}

type VoteResponse struct {
	ID uuid.UUID
}
