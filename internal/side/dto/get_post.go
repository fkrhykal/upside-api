package dto

import (
	"github.com/fkrhykal/upside-api/internal/shared/pagination"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type GetPostsResponse struct {
	Posts    Posts                           `json:"posts"`
	Metadata *pagination.CursorBasedMetadata `json:"metadata"`
}

type Post struct {
	ID        ulid.ULID `json:"id"`
	Body      string    `json:"body"`
	CreatedAt int64     `json:"createdAt"`
	UpdatedAt int64     `json:"updatedAt"`
	Author    *Author   `json:"author"`
	Side      *Side     `json:"side"`
}

type Posts []*Post

var EmptyPosts = make(Posts, 0)

type Author struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func MapPosts(posts entity.Posts) Posts {
	postsDto := make(Posts, len(posts))
	for i, post := range posts {
		postsDto[i] = MapPost(post)
	}
	return postsDto
}

func MapPost(post *entity.Post) *Post {
	return &Post{
		ID:        post.ID,
		Body:      post.Body,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Author: &Author{
			ID:       post.Author.ID,
			Username: post.Author.Username,
		},
		Side: &Side{
			ID:   post.Side.ID,
			Nick: post.Side.Nick,
			Name: post.Side.Name,
		},
	}
}
