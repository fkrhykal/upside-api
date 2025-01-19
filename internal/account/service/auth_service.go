package service

import (
	"context"

	"github.com/fkrhykal/upside-api/internal/account/dto"
)

type AuthService interface {
	SignUp(ctx context.Context, request *dto.SignUpRequest) (*dto.SignUpResponse, error)
	SignIn(ctx context.Context, request *dto.SignInRequest) (*dto.SignInResponse, error)
}
