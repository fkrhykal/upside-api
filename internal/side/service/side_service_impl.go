package service

import (
	"context"
	"time"

	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/validation"
	"github.com/fkrhykal/upside-api/internal/side/dto"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/fkrhykal/upside-api/internal/side/repository"
	"github.com/google/uuid"
)

type SideServiceImpl[T any] struct {
	logger               log.Logger
	validator            validation.Validator
	ctxManager           db.CtxManager[T]
	sideRepository       repository.SideRepository[T]
	membershipRepository repository.MembershipRepository[T]
}

func NewSideServiceImpl[T any](
	logger log.Logger,
	validator validation.Validator,
	ctxManager db.CtxManager[T],
	sideRepository repository.SideRepository[T],
	membershipRepository repository.MembershipRepository[T],
) SideService {
	return &SideServiceImpl[T]{
		logger:               logger,
		validator:            validator,
		ctxManager:           ctxManager,
		sideRepository:       sideRepository,
		membershipRepository: membershipRepository,
	}
}

func (s *SideServiceImpl[T]) CreateSide(ctx context.Context, req *dto.CreateSideRequest) (*dto.CreateSideResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	tx, err := s.ctxManager.NewTxContext(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	side := &entity.Side{
		ID:          uuid.New(),
		Nick:        req.Nick,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   uint64(time.Now().UnixMilli()),
	}
	if err := s.sideRepository.Save(tx, side); err != nil {
		return nil, err
	}

	membership := &entity.Membership{
		ID:     uuid.New(),
		Member: req.FounderID,
		Side:   side.ID,
		Role:   entity.FOUNDER,
	}
	if err := s.membershipRepository.Save(tx, membership); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &dto.CreateSideResponse{ID: side.ID}, nil
}
