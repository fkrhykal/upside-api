package dto

import "github.com/oklog/ulid/v2"

type RevokeVoteRequest struct {
	PostID ulid.ULID
}
