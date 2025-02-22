package repository

import (
	"database/sql"
	"errors"

	c "github.com/fkrhykal/upside-api/internal/shared/collection"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/oklog/ulid/v2"
)

type PgVoteRepository struct {
	logger log.Logger
}

func (pv *PgVoteRepository) Save(ctx db.DBContext[db.SqlExecutor], vote *entity.Vote) error {
	query := `INSERT INTO votes(id, voter_id, post_id, kind) VALUES($1, $2, $3, $4)`
	_, err := ctx.Executor().ExecContext(ctx, query, vote.ID, vote.Voter.ID, vote.Post.ID.String(), vote.Kind)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (pv *PgVoteRepository) Update(ctx db.DBContext[db.SqlExecutor], vote *entity.Vote) error {
	query := `UPDATE votes SET kind = $1 WHERE id = $2`
	_, err := ctx.Executor().ExecContext(ctx, query, vote.Kind, vote.ID)
	if err != nil {
		return err
	}
	return nil
}

func (pv *PgVoteRepository) FindByVoterIDAndPostID(ctx db.DBContext[db.SqlExecutor], voterID uuid.UUID, postID ulid.ULID) (*entity.Vote, error) {
	vote := &entity.Vote{
		Voter: &entity.Voter{},
		Post:  &entity.Post{},
	}
	var id string
	query := `SELECT id, voter_id, post_id, kind FROM votes WHERE voter_id = $1 AND post_id = $2`
	err := ctx.Executor().
		QueryRowContext(ctx, query, voterID, postID.String()).
		Scan(&vote.ID, &vote.Voter.ID, &id, &vote.Kind)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	postULID, err := ulid.Parse(id)
	if err != nil {
		return nil, err
	}
	vote.Post.ID = postULID
	return vote, nil
}

func (pv *PgVoteRepository) SumVoteKindGroupByPostIDs(ctx db.DBContext[db.SqlExecutor], postIDs []ulid.ULID) (map[ulid.ULID]int, error) {
	query := `SELECT SUM(v.kind), p.id FROM votes AS v JOIN posts AS p ON v.post_id = p.id WHERE p.id = ANY($1) GROUP BY p.id`
	row, err := ctx.Executor().QueryContext(ctx, query, pq.Array(c.Map(postIDs, func(postID ulid.ULID) string { return postID.String() })))
	if err != nil {
		return nil, err
	}
	registry := make(map[ulid.ULID]int)
	for row.Next() {
		var id string
		var vote_score int
		if err := row.Scan(&vote_score, &id); err != nil {
			return nil, err
		}
		postID, err := ulid.Parse(id)
		if err != nil {
			return nil, err
		}
		registry[postID] = vote_score
	}
	return registry, nil
}

func (pv *PgVoteRepository) FindManyByVoterIDAndPostIDs(ctx db.DBContext[db.SqlExecutor], voterID uuid.UUID, postIDs []ulid.ULID) ([]*entity.Vote, error) {
	query := `SELECT id, voter_id, post_id, kind FROM votes WHERE post_id = ANY($1) AND voter_id = $2`
	row, err := ctx.Executor().QueryContext(ctx, query, pq.Array(c.Map(postIDs, func(id ulid.ULID) string { return id.String() })), voterID)
	if err != nil {
		return nil, err
	}
	var votes []*entity.Vote
	for row.Next() {
		var postID string
		vote := &entity.Vote{
			Voter: &entity.Voter{},
			Post:  &entity.Post{},
		}
		if err := row.Scan(&vote.ID, &vote.Voter.ID, &postID, &vote.Kind); err != nil {
			return nil, err
		}
		postULID, err := ulid.Parse(postID)
		if err != nil {
			return nil, err
		}
		vote.Post.ID = postULID
		votes = append(votes, vote)
		pv.logger.Debugf("%+s", vote.ID)
	}
	if err := row.Err(); err != nil {
		return nil, err
	}
	return votes, err
}

func (pv *PgVoteRepository) DeleteVote(ctx db.DBContext[db.SqlExecutor], voteID uuid.UUID) error {
	query := `DELETE FROM votes WHERE id = $1`
	_, err := ctx.Executor().ExecContext(ctx, query, voteID)
	if err != nil {
		return err
	}
	return nil
}

func NewPgVoteRepository(logger log.Logger) VoteRepository[db.SqlExecutor] {
	return &PgVoteRepository{logger: logger}
}
