package service

import (
	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/exception"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/side/dto"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/fkrhykal/upside-api/internal/side/repository"
	"github.com/google/uuid"
)

type VoteServiceImpl[T any] struct {
	logger               log.Logger
	ctxManager           db.CtxManager[T]
	sideRepository       repository.SideRepository[T]
	membershipRepository repository.MembershipRepository[T]
	postRepository       repository.PostRepository[T]
	voteRepository       repository.VoteRepository[T]
}

func NewVoteServiceImpl[T any](
	logger log.Logger,
	ctxManager db.CtxManager[T],
	sideRepository repository.SideRepository[T],
	membershipRepository repository.MembershipRepository[T],
	postRepository repository.PostRepository[T],
	voteRepository repository.VoteRepository[T],
) VoteService {
	return &VoteServiceImpl[T]{
		logger:               logger,
		ctxManager:           ctxManager,
		sideRepository:       sideRepository,
		membershipRepository: membershipRepository,
		postRepository:       postRepository,
		voteRepository:       voteRepository,
	}
}

func (vs *VoteServiceImpl[T]) Vote(ctx *auth.AuthContext, req *dto.VoteRequest) (*dto.VoteResponse, error) {
	dbCtx := vs.ctxManager.NewDBContext(ctx)
	post, err := vs.postRepository.FindByID(dbCtx, req.PostID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, &exception.NotFoundError{}
	}
	membership, err := vs.membershipRepository.FindBySideIDAndMemberID(dbCtx, ctx.Credential.ID, post.Side.ID)
	if err != nil {
		return nil, err
	}
	if membership == nil {
		return nil, exception.ErrAuthorization
	}
	vote, err := vs.voteRepository.FindByVoterIDAndPostID(dbCtx, ctx.Credential.ID, post.ID)
	if err != nil {
		return nil, err
	}
	if vote == nil {
		vote = &entity.Vote{
			ID:    uuid.New(),
			Voter: &entity.Voter{ID: ctx.Credential.ID},
			Post:  &entity.Post{ID: post.ID},
			Kind:  req.VoteKind,
		}
		if err := vs.voteRepository.Save(dbCtx, vote); err != nil {
			return nil, err
		}
	}
	if vote.Kind == req.VoteKind {
		return &dto.VoteResponse{ID: vote.ID}, nil
	}
	vote.Kind = req.VoteKind
	if err := vs.voteRepository.Update(dbCtx, vote); err != nil {
		return nil, err
	}
	return &dto.VoteResponse{ID: vote.ID}, nil
}

func (vs *VoteServiceImpl[T]) RevokeVote(ctx *auth.AuthContext, req *dto.RevokeVoteRequest) error {
	dbCtx := vs.ctxManager.NewDBContext(ctx)
	post, err := vs.postRepository.FindByID(dbCtx, req.PostID)
	if err != nil {
		return err
	}
	if post == nil {
		return &exception.NotFoundError{ResourceName: "post"}
	}
	membership, err := vs.membershipRepository.FindBySideIDAndMemberID(dbCtx, ctx.Credential.ID, post.Side.ID)
	if err != nil {
		return err
	}
	if membership == nil {
		return exception.ErrAuthorization
	}
	vote, err := vs.voteRepository.FindByVoterIDAndPostID(dbCtx, ctx.Credential.ID, post.ID)
	if err != nil {
		return err
	}
	if vote == nil {
		return &exception.NotFoundError{ResourceName: "vote"}
	}
	if err := vs.voteRepository.DeleteVote(dbCtx, vote.ID); err != nil {
		return err
	}
	return nil
}
