package service

import (
	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/collection"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/exception"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/pagination"
	"github.com/fkrhykal/upside-api/internal/shared/validation"
	"github.com/fkrhykal/upside-api/internal/side/dto"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/fkrhykal/upside-api/internal/side/repository"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type PostServiceImpl[T any] struct {
	logger               log.Logger
	validator            validation.Validator
	ctxManager           db.CtxManager[T]
	sideRepository       repository.SideRepository[T]
	membershipRepository repository.MembershipRepository[T]
	postRepository       repository.PostRepository[T]
	voteRepository       repository.VoteRepository[T]
}

func NewPostServiceImpl[T any](
	logger log.Logger,
	validator validation.Validator,
	ctxManager db.CtxManager[T],
	sideRepository repository.SideRepository[T],
	membershipRepository repository.MembershipRepository[T],
	postRepository repository.PostRepository[T],
	voteRepository repository.VoteRepository[T],
) PostService {
	return &PostServiceImpl[T]{
		logger:               logger,
		validator:            validator,
		ctxManager:           ctxManager,
		sideRepository:       sideRepository,
		membershipRepository: membershipRepository,
		postRepository:       postRepository,
		voteRepository:       voteRepository,
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

	var sideIDs uuid.UUIDs
	postIDs := make([]entity.PostID, len(m.Data))

	for i, post := range m.Data {
		postIDs[i] = post.ID
		sideIDs = append(sideIDs, post.Side.ID)
	}

	scoreRegistry, err := ps.voteRepository.SumVoteKindGroupByPostIDs(dbCtx, postIDs)
	if err != nil {
		return nil, err
	}

	voteRegistry := make(map[entity.PostID]*entity.Vote)
	membershipDetailRegistry := make(map[entity.PostID]*entity.Membership)

	if ctx.Authenticated() {
		votes, err := ps.voteRepository.FindManyByVoterIDAndPostIDs(dbCtx, ctx.Credential.ID, postIDs)
		if err != nil {
			return nil, err
		}
		for _, vote := range votes {
			voteRegistry[vote.Post.ID] = vote
		}
		memberships, err := ps.membershipRepository.FindManyBySideIDsAndMemberID(dbCtx, sideIDs, ctx.Credential.ID)
		if err != nil {
			return nil, err
		}
		membershipRegistry := make(map[uuid.UUID]*entity.Membership, len(memberships))
		for _, membership := range memberships {
			membershipRegistry[membership.Side] = membership
		}
		for _, post := range m.Data {
			if membership, ok := membershipRegistry[post.Side.ID]; ok {
				membershipDetailRegistry[post.ID] = membership
			}
		}
	}

	return &dto.GetPostsResponse{
		Posts:    dto.MapPosts(m.Data, scoreRegistry, voteRegistry, membershipDetailRegistry),
		Metadata: m.Metadata,
	}, nil
}

func (ps *PostServiceImpl[T]) GetSubscribedPosts(ctx *auth.AuthContext, cursor pagination.ULIDCursor) (*dto.GetPostsResponse, error) {
	if !ctx.Authenticated() {
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

	sideIDs := collection.Map(memberships, func(m *entity.Membership) uuid.UUID { return m.Side })

	m, err := ps.postRepository.FindManyWithULIDCursorInSides(dbCtx, cursor, sideIDs)
	if err != nil {
		return nil, err
	}

	postIDs := collection.Map(m.Data, func(post *entity.Post) ulid.ULID {
		return post.ID
	})

	membershipDetailRegistry := make(map[entity.PostID]*entity.Membership)

	membershipRegistry := make(map[uuid.UUID]*entity.Membership, len(memberships))
	for _, membership := range memberships {
		membershipRegistry[membership.Side] = membership
	}
	for _, post := range m.Data {
		if membership, ok := membershipRegistry[post.Side.ID]; ok {
			membershipDetailRegistry[post.ID] = membership
		}
	}

	scoreRegistry, err := ps.voteRepository.SumVoteKindGroupByPostIDs(dbCtx, postIDs)
	if err != nil {
		return nil, err
	}

	voteRegistry := make(map[entity.PostID]*entity.Vote)

	votes, err := ps.voteRepository.FindManyByVoterIDAndPostIDs(dbCtx, ctx.Credential.ID, postIDs)
	if err != nil {
		return nil, err
	}
	for _, vote := range votes {
		voteRegistry[vote.Post.ID] = vote
	}

	return &dto.GetPostsResponse{
		Posts:    dto.MapPosts(m.Data, scoreRegistry, voteRegistry, membershipDetailRegistry),
		Metadata: m.Metadata,
	}, nil
}
