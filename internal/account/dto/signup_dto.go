package dto

import "github.com/google/uuid"

type SignUpRequest struct {
	Username string `json:"username" validate:"required,min=4,max=16,alphanum,ascii" name:"username"`
	Password string `json:"password" validate:"required,min=8,max=128,password" name:"password"`
} //@name SignUpRequest

type SignUpResponse struct {
	ID uuid.UUID `json:"id" example:"b8ae9666-23b8-4ffb-8cf1-f7df0d02130a"`
} //@name SignUpResponse
