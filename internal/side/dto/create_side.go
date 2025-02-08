package dto

import "github.com/google/uuid"

type CreateSideRequest struct {
	Nick        string `json:"nick" validate:"required,alphanum,min=4,max=24"`
	Name        string `json:"name" validate:"required,min=4,max=32"`
	Description string `json:"description" validate:"required,max=1000"`
}

type CreateSideResponse struct {
	ID uuid.UUID
}
