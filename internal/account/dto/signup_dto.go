package dto

import "github.com/google/uuid"

type SignUpRequest struct {
	Username string `json:"username" validate:"required" name:"username"`
	Password string `json:"password" validate:"required" name:"password"`
}

type SignUpResponse struct {
	ID uuid.UUID `json:"id"`
}
