package repository

import (
	"database/sql"
	"errors"

	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/google/uuid"
)

type PgSideRepository struct {
	logger log.Logger
}

func NewPgSideRepository(logger log.Logger) *PgSideRepository {
	return &PgSideRepository{
		logger: logger,
	}
}

func (p *PgSideRepository) Save(ctx db.DBContext[db.SqlExecutor], side *entity.Side) error {
	query := "INSERT INTO sides(id, nick, name, description, created_at) VALUES($1,$2,$3,$4,$5)"

	_, err := ctx.Executor().
		ExecContext(ctx, query, side.ID, side.Nick, side.Name, side.Description, side.CreatedAt)
	if err != nil {
		return err
	}
	p.logger.Debugf(
		"Saved data to table sides(id=%s, nick=%s, name=%s, description=%s, created_at=%d)",
		side.ID,
		side.Nick,
		side.Name,
		side.Description,
		side.CreatedAt,
	)
	return nil
}

func (p *PgSideRepository) FindById(ctx db.DBContext[db.SqlExecutor], id uuid.UUID) (*entity.Side, error) {
	query := "SELECT id, nick, name, description, created_at FROM sides WHERE id = $1"

	side := new(entity.Side)

	err := ctx.Executor().
		QueryRowContext(ctx, query, id).
		Scan(&side.ID, &side.Nick, &side.Name, &side.Description, &side.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	p.logger.Debugf(
		"Retrieved data from table sides(id=%s, nick=%s, name=%s, description=%s, created_at=%d)",
		side.ID,
		side.Nick,
		side.Name,
		side.Description,
		side.CreatedAt,
	)
	return side, nil
}
