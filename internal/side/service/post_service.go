package service

import (
	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/pagination"
	"github.com/fkrhykal/upside-api/internal/side/dto"
)

type PostService interface {
	CreatePost(ctx *auth.AuthContext, req *dto.CreatePostRequest) (*dto.CreatePostResponse, error)
	GetLatestPosts(ctx *auth.AuthContext, cursor pagination.ULIDCursor) (*dto.GetPostsResponse, error)
	GetSubscribedPosts(ctx *auth.AuthContext, cursor pagination.ULIDCursor) (*dto.GetPostsResponse, error)
}
