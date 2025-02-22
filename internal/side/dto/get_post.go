package dto

import (
	"github.com/fkrhykal/upside-api/internal/shared/collection"
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
	ID              ulid.ULID        `json:"id"`
	Body            string           `json:"body"`
	CreatedAt       int64            `json:"createdAt"`
	UpdatedAt       int64            `json:"updatedAt"`
	Author          *Author          `json:"author"`
	Side            *Side            `json:"side"`
	Score           int              `json:"score"`
	CurrentUserVote *CurrentUserVote `json:"currentUserVote"`
}

type CurrentUserVote struct {
	ID   uuid.UUID       `json:"id"`
	Kind entity.VoteKind `json:"kind"`
}

type Posts []*Post

var EmptyPosts = make(Posts, 0)

type Author struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func MapPosts(posts entity.Posts, scoreRegistry map[ulid.ULID]int, voteRegistry map[ulid.ULID]*entity.Vote, membershipDetailRegistry map[ulid.ULID]*entity.Membership) Posts {
	return collection.Map(posts, func(post *entity.Post) *Post {

		var currentUserVote *CurrentUserVote

		if vote, ok := voteRegistry[post.ID]; ok {
			currentUserVote = &CurrentUserVote{
				ID:   vote.ID,
				Kind: vote.Kind,
			}
		}

		var membershipDetail *MembershipDetail

		if membership, ok := membershipDetailRegistry[post.ID]; ok {
			membershipDetail = &MembershipDetail{
				ID:   membership.ID,
				Role: membership.Role.String(),
			}
		}

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
				ID:               post.Side.ID,
				Nick:             post.Side.Nick,
				Name:             post.Side.Name,
				MembershipDetail: membershipDetail,
			},
			Score:           scoreRegistry[post.ID],
			CurrentUserVote: currentUserVote,
		}
	})
}
