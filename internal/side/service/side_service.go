package service

import (
	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/pagination"
	"github.com/fkrhykal/upside-api/internal/side/dto"
)

type SideService interface {
	CreateSide(ctx *auth.AuthContext, req *dto.CreateSideRequest) (*dto.CreateSideResponse, error)
	GetJoinedSides(ctx *auth.AuthContext, page *pagination.OffsetBased) (dto.Sides, error)
	GetPopularSides(ctx *auth.AuthContext, page *pagination.OffsetBased) (dto.Sides, error)
	GetSides(ctx *auth.AuthContext, page *pagination.OffsetBased) (dto.Sides, error)
}
