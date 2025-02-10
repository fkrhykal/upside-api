package repository

import (
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/side/entity"
)

// id
// body
// created_at
// updated_at
// author_id
// side_id
type PgPostRepository struct {
	logger log.Logger
}

func NewPgPostRepository(logger log.Logger) PostRepository[db.SqlExecutor] {
	return &PgPostRepository{
		logger: logger,
	}
}

func (pr *PgPostRepository) Save(ctx db.DBContext[db.SqlExecutor], post *entity.Post) error {
	query := `INSERT INTO posts (id, body, created_at, updated_at, author_id, side_id) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := ctx.
		Executor().
		ExecContext(ctx, query, post.ID.String(), post.Body, post.CreatedAt, post.UpdatedAt, post.Author.ID, post.Side.ID)
	return err
}
