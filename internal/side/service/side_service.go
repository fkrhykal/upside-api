package service

import (
	"context"

	"github.com/fkrhykal/upside-api/internal/side/dto"
)

type SideService interface {
	CreateSide(ctx context.Context, req *dto.CreateSideRequest) (*dto.CreateSideResponse, error)
}
