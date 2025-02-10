package service

import (
	"time"

	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/exception"
	"github.com/fkrhykal/upside-api/internal/shared/validation"
	"github.com/fkrhykal/upside-api/internal/side/dto"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/fkrhykal/upside-api/internal/side/repository"
	"github.com/oklog/ulid/v2"
)

type PostServiceImpl[T any] struct {
	validator            validation.Validator
	ctxManager           db.CtxManager[T]
	sideRepository       repository.SideRepository[T]
	membershipRepository repository.MembershipRepository[T]
	postRepository       repository.PostRepository[T]
}

func NewPostServiceImpl[T any](
	validator validation.Validator,
	ctxManager db.CtxManager[T],
	sideRepository repository.SideRepository[T],
	membershipRepository repository.MembershipRepository[T],
	postRepository repository.PostRepository[T],
) PostService {
	return &PostServiceImpl[T]{
		validator:            validator,
		ctxManager:           ctxManager,
		sideRepository:       sideRepository,
		membershipRepository: membershipRepository,
		postRepository:       postRepository,
	}
}

func (ps *PostServiceImpl[T]) CreatePost(ctx *auth.AuthContext, req *dto.CreatePostRequest) (*dto.CreatePostResponse, error) {
	if err := ps.validator.Validate(req); err != nil {
		return nil, err
	}
	dbCtx := ps.ctxManager.NewDBContext(ctx)
	side, err := ps.sideRepository.FindById(dbCtx, req.SideID)
	if err != nil {
		return nil, err
	}
	if side == nil {
		return nil, exception.ErrSideNotFound
	}

	membership, err := ps.membershipRepository.FindBySideIDAndMemberID(dbCtx, ctx.Credential.ID, side.ID)
	if err != nil {
		return nil, err
	}
	if membership == nil {
		return nil, exception.ErrAuthorization
	}

	post := &entity.Post{
		ID:        ulid.Make(),
		Body:      req.Body,
		CreatedAt: time.Now().UnixMilli(),
		Author:    &entity.Author{ID: ctx.Credential.ID},
		Side:      side,
	}

	if err := ps.postRepository.Save(dbCtx, post); err != nil {
		return nil, err
	}

	return &dto.CreatePostResponse{ID: post.ID}, nil
}
