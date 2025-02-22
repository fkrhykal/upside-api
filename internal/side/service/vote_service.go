package service

import (
	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/side/dto"
)

type VoteService interface {
	Vote(ctx *auth.AuthContext, req *dto.VoteRequest) (*dto.VoteResponse, error)
	RevokeVote(ctx *auth.AuthContext, req *dto.RevokeVoteRequest) error
}
