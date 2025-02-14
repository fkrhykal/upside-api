package service

import (
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

type PostServiceImpl[T any] struct {
	logger               log.Logger
	validator            validation.Validator
	ctxManager           db.CtxManager[T]
	sideRepository       repository.SideRepository[T]
	membershipRepository repository.MembershipRepository[T]
	postRepository       repository.PostRepository[T]
}

func NewPostServiceImpl[T any](
	logger log.Logger,
	validator validation.Validator,
	ctxManager db.CtxManager[T],
	sideRepository repository.SideRepository[T],
	membershipRepository repository.MembershipRepository[T],
	postRepository repository.PostRepository[T],
) PostService {
	return &PostServiceImpl[T]{
		logger:               logger,
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

	post := entity.CreatePost(req.Body, ctx.Credential.ID, side.ID)

	if err := ps.postRepository.Save(dbCtx, post); err != nil {
		return nil, err
	}

	return &dto.CreatePostResponse{ID: post.ID}, nil
}

func (ps *PostServiceImpl[T]) GetLatestPosts(ctx *auth.AuthContext, cursor pagination.ULIDCursor) (*dto.GetPostsResponse, error) {
	dbCtx := ps.ctxManager.NewDBContext(ctx)

	m, err := ps.postRepository.FindManyWithULIDCursor(dbCtx, cursor)
	if err != nil {
		return nil, err
	}

	return &dto.GetPostsResponse{
		Posts:    dto.MapPosts(m.Data),
		Metadata: m.Metadata,
	}, nil
}

func (ps *PostServiceImpl[T]) GetSubscribedPosts(ctx *auth.AuthContext, cursor pagination.ULIDCursor) (*dto.GetPostsResponse, error) {
	if !ctx.Authenticated() {
		ps.logger.Debugf("authenticated user: %s", ctx.Authenticated())
		return &dto.GetPostsResponse{Posts: dto.EmptyPosts, Metadata: &pagination.CursorBasedMetadata{}}, nil
	}

	dbCtx := ps.ctxManager.NewDBContext(ctx)
	memberships, err := ps.membershipRepository.FindManyByMemberID(dbCtx, ctx.Credential.ID)

	if err != nil {
		return nil, err
	}
	if len(memberships) == 0 {
		return &dto.GetPostsResponse{Posts: dto.EmptyPosts, Metadata: &pagination.CursorBasedMetadata{}}, nil
	}

	sideIDs := make([]uuid.UUID, len(memberships))
	for i, membership := range memberships {
		sideIDs[i] = membership.Side
	}

	m, err := ps.postRepository.FindManyWithULIDCursorInSides(dbCtx, cursor, sideIDs)
	if err != nil {
		return nil, err
	}

	return &dto.GetPostsResponse{
		Posts:    dto.MapPosts(m.Data),
		Metadata: m.Metadata,
	}, nil
}
