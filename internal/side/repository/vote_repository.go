package repository

import (
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type VoteRepository[T any] interface {
	Save(ctx db.DBContext[T], vote *entity.Vote) error
	Update(ctx db.DBContext[T], vote *entity.Vote) error
	FindByVoterIDAndPostID(ctx db.DBContext[T], voterID uuid.UUID, postID ulid.ULID) (*entity.Vote, error)
	DeleteVote(ctx db.DBContext[T], voteID uuid.UUID) error
	SumVoteKindGroupByPostIDs(ctx db.DBContext[T], postIDs []ulid.ULID) (map[ulid.ULID]int, error)
	FindManyByVoterIDAndPostIDs(ctx db.DBContext[T], voterID uuid.UUID, postIDs []ulid.ULID) ([]*entity.Vote, error)
}
