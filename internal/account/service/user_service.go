package service

import (
	"context"

	"github.com/fkrhykal/upside-api/internal/account/dto"
	"github.com/google/uuid"
)

type UserService interface {
	GetUserDetail(ctx context.Context, id uuid.UUID) (*dto.UserDetail, error)
}
