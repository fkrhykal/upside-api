package dto

import "github.com/google/uuid"

type CreateSideRequest struct {
	Nick        string    `json:"nick" validate:"required,alphanum"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	FounderID   uuid.UUID `json:"_"`
}

type CreateSideResponse struct {
	ID uuid.UUID
}
