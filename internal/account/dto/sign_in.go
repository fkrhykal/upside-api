package dto

import "github.com/fkrhykal/upside-api/internal/shared/auth"

type SignInRequest struct {
	Username string `json:"username" validate:"required,min=4,max=16,alphanum,ascii" name:"username"`
	Password string `json:"password" validate:"required,min=8,max=128,password" name:"password"`
} //@name SignInRequest

type SignInResponse struct {
	Token auth.Token  `json:"token"`
	User  *UserDetail `json:"user"`
} //@name SignInResponse
