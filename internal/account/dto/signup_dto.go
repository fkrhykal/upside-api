package dto

import "github.com/google/uuid"

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	ID uuid.UUID `json:"id"`
}
