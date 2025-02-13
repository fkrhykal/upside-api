package service

import (
	"errors"

	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/exception"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/pagination"
	"github.com/fkrhykal/upside-api/internal/shared/validation"
	"github.com/fkrhykal/upside-api/internal/side/dto"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/fkrhykal/upside-api/internal/side/repository"
	"github.com/oklog/ulid/v2"
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

	post := entity.NewPost(req.Body, ctx.Credential.ID, side.ID)

	if err := ps.postRepository.Save(dbCtx, post); err != nil {
		return nil, err
	}

	return &dto.CreatePostResponse{ID: post.ID}, nil
}

func (ps *PostServiceImpl[T]) GetLatestPosts(ctx *auth.AuthContext, cursor pagination.ULIDCursor) (*dto.GetPostsResponse, error) {
	dbCtx := ps.ctxManager.NewDBContext(ctx)
	var nextCursor *pagination.NextULIDCursor
	var prevCursor *pagination.PrevULIDCursor

	if cursor.ID() == nil {
		posts, err := ps.postRepository.FindManyWithULIDCursor(dbCtx, cursor)
		if err != nil {
			return nil, err
		}
		if len(posts) <= cursor.Limit() {
			nextCursor = pagination.LimitNextULIDCursor(cursor.Limit())
		} else {
			var postID *ulid.ULID

			if len(posts) < cursor.Limit()+1 {
				postID = &posts.Last().ID
			} else {
				postID = &posts.At(posts.LastIndex() - 2).ID
			}

			nextCursor, err = pagination.NewNextULIDCursor(postID, cursor.Limit())
			if err != nil {
				return nil, err
			}
		}

		prevCursor = pagination.LimitPrevULIDCursor(cursor.Limit())

		return &dto.GetPostsResponse{
			Posts:    dto.FromEntitiesToPosts(posts.Slice(0, cursor.Limit())),
			Metadata: pagination.ULIDCursorMetadata(prevCursor, nextCursor),
		}, nil
	}

	_, isNextCursor := cursor.(*pagination.NextULIDCursor)
	if isNextCursor {
		posts, err := ps.postRepository.FindManyWithULIDCursor(dbCtx, cursor)
		if err != nil {
			return nil, err
		}
		if len(posts) < cursor.Limit()+2 {
			nextCursor = pagination.LimitNextULIDCursor(cursor.Limit())
		} else {
			nextCursor, err = pagination.NewNextULIDCursor(&posts.Penultimate().ID, cursor.Limit())
			if err != nil {
				return nil, err
			}
		}
		prevCursor, err = pagination.NewPrevULIDCursor(cursor.ID(), cursor.Limit())
		if err != nil {
			return nil, err
		}
		return &dto.GetPostsResponse{
			Posts:    dto.FromEntitiesToPosts(posts.Slice(1, cursor.Limit()+1)),
			Metadata: pagination.ULIDCursorMetadata(prevCursor, nextCursor),
		}, nil
	}

	_, isPrevCursor := cursor.(*pagination.PrevULIDCursor)
	if isPrevCursor {
		posts, err := ps.postRepository.FindManyWithULIDCursor(dbCtx, cursor)
		if err != nil {
			return nil, err
		}
		if len(posts) < cursor.Limit()+2 {
			prevCursor = pagination.LimitPrevULIDCursor(cursor.Limit())
		} else {
			prevCursor, err = pagination.NewPrevULIDCursor(&posts.Second().ID, cursor.Limit())
			if err != nil {
				return nil, err
			}
		}

		nextCursor, err = pagination.NewNextULIDCursor(cursor.ID(), cursor.Limit())
		if err != nil {
			return nil, err
		}
		var prevPosts entity.Posts
		if len(posts) < cursor.Limit()+2 {
			prevPosts = posts
		} else {
			prevPosts = posts[2:]
		}
		return &dto.GetPostsResponse{
			Posts:    dto.FromEntitiesToPosts(prevPosts),
			Metadata: pagination.ULIDCursorMetadata(prevCursor, nextCursor),
		}, nil
	}

	return nil, errors.New("cursor is nil")
}
