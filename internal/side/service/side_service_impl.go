package service

import (
	"time"

	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/exception"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/pagination"
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

func (s *SideServiceImpl[T]) GetSides(ctx *auth.AuthContext, page *pagination.OffsetBased) (dto.Sides, error) {
	dbCtx := s.ctxManager.NewDBContext(ctx)

	sides, err := s.sideRepository.FindManyWithOffsetAndLimit(dbCtx, page.Offset(), page.Limit)
	if err != nil {
		return dto.EmptySides, err
	}

	if !ctx.Authenticated() {
		sidesDto := make(dto.Sides, len(sides))
		for i, side := range sides {
			sidesDto[i] = &dto.Side{ID: side.ID, Nick: side.Nick, Name: side.Description}
		}
		return sidesDto, nil
	}
	sideIDs := make(uuid.UUIDs, len(sides))

	for i, side := range sides {
		sideIDs[i] = side.ID
	}

	memberships, err := s.membershipRepository.FindManyBySideIDsAndMemberID(dbCtx, sideIDs, ctx.Credential.ID)
	if err != nil {
		return dto.EmptySides, err
	}

	membershipRegistry := make(map[uuid.UUID]*entity.Membership, len(memberships))

	for _, membership := range memberships {
		membershipRegistry[membership.Side] = membership
	}

	sidesDto := make(dto.Sides, len(sides))

	for i, side := range sides {
		sideDto := &dto.Side{ID: side.ID, Nick: side.Nick, Name: side.Description}
		if membership, ok := membershipRegistry[side.ID]; ok {
			sideDto.MembershipDetail = &dto.MembershipDetail{
				ID:   membership.ID,
				Role: membership.Role.String(),
			}
		}
		sidesDto[i] = sideDto
	}
	return sidesDto, nil
}

func (s *SideServiceImpl[T]) GetPopularSides(ctx *auth.AuthContext, page *pagination.OffsetBased) (dto.Sides, error) {
	dbCtx := s.ctxManager.NewDBContext(ctx)

	sides, err := s.sideRepository.FindOffsetLimitedWithLargestMemberships(dbCtx, page.Offset(), page.Limit)
	if err != nil {
		return make(dto.Sides, 0), nil
	}
	if !ctx.Authenticated() {
		sidesDto := make(dto.Sides, len(sides))
		for i, side := range sides {
			sidesDto[i] = &dto.Side{ID: side.ID, Nick: side.Nick, Name: side.Description}
		}
		return sidesDto, nil
	}
	sideIDs := make(uuid.UUIDs, len(sides))

	for i, side := range sides {
		sideIDs[i] = side.ID
	}

	memberships, err := s.membershipRepository.FindManyBySideIDsAndMemberID(dbCtx, sideIDs, ctx.Credential.ID)
	if err != nil {
		return make(dto.Sides, 0), err
	}

	membershipRegistry := make(map[uuid.UUID]*entity.Membership, len(memberships))

	for _, membership := range memberships {
		membershipRegistry[membership.Side] = membership
	}

	sidesDto := make(dto.Sides, len(sides))

	for i, side := range sides {
		sideDto := &dto.Side{ID: side.ID, Nick: side.Nick, Name: side.Description}
		if membership, ok := membershipRegistry[side.ID]; ok {
			sideDto.MembershipDetail = &dto.MembershipDetail{
				ID:   membership.ID,
				Role: membership.Role.String(),
			}
		}
		sidesDto[i] = sideDto
	}
	return sidesDto, nil
}

func (s *SideServiceImpl[T]) GetJoinedSides(ctx *auth.AuthContext, page *pagination.OffsetBased) (dto.Sides, error) {
	if !ctx.Authenticated() {
		return dto.EmptySides, nil
	}

	dbCtx := s.ctxManager.NewDBContext(ctx)

	memberships, err := s.membershipRepository.FindOffsetLimitedByMemberID(dbCtx, ctx.Credential.ID, page.Offset(), page.Limit)
	if err != nil {
		return dto.EmptySides, err
	}

	sideIDs := make([]uuid.UUID, len(memberships))
	membershipRegistry := make(map[uuid.UUID]*entity.Membership, len(memberships))

	for i, membership := range memberships {
		sideIDs[i] = membership.Side
		membershipRegistry[membership.Side] = membership
	}

	sides, err := s.sideRepository.FindManyIn(dbCtx, sideIDs)
	if err != nil {
		return dto.EmptySides, err
	}

	sidesDto := make(dto.Sides, len(sides))

	for i, side := range sides {
		membership := membershipRegistry[side.ID]
		membershipDetail := &dto.MembershipDetail{
			ID:   membership.ID,
			Role: string(membership.Role),
		}
		sidesDto[i] = &dto.Side{
			ID:               side.ID,
			Nick:             side.Nick,
			Name:             side.Name,
			MembershipDetail: membershipDetail,
		}
	}

	return sidesDto, nil
}

func (s *SideServiceImpl[T]) CreateSide(ctx *auth.AuthContext, req *dto.CreateSideRequest) (*dto.CreateSideResponse, error) {
	if !ctx.Authenticated() {
		return nil, exception.ErrAuthentication
	}

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
		Member: ctx.Credential.ID,
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
