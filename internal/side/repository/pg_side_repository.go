package repository

import (
	"database/sql"
	"errors"

	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PgSideRepository struct {
	logger log.Logger
}

func NewPgSideRepository(logger log.Logger) SideRepository[db.SqlExecutor] {
	return &PgSideRepository{
		logger: logger,
	}
}

func (p *PgSideRepository) Save(ctx db.DBContext[db.SqlExecutor], side *entity.Side) error {
	query := `INSERT INTO sides(id, nick, name, description, created_at) VALUES($1,$2,$3,$4,$5)`

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
	query := `SELECT id, nick, name, description, created_at FROM sides WHERE id = $1`

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

func (p *PgSideRepository) FindManyIn(ctx db.DBContext[db.SqlExecutor], ids []uuid.UUID) ([]*entity.Side, error) {
	query := `SELECT id, nick, name, description, created_at FROM sides WHERE id = ANY($1)`

	rows, err := ctx.Executor().QueryContext(ctx, query, pq.Array(ids))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var sides []*entity.Side

	for rows.Next() {
		side := new(entity.Side)
		if err := rows.Scan(&side.ID, &side.Nick, &side.Name, &side.Description, &side.CreatedAt); err != nil {
			continue
		}
		sides = append(sides, side)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sides, nil
}

func (p *PgSideRepository) FindLimitedWithLargestMemberships(ctx db.DBContext[db.SqlExecutor], limit int) ([]*entity.Side, error) {
	query := `SELECT s.id, s.nick, s.name, s.description, s.created_at, 
    (SELECT COUNT(*) FROM memberships m WHERE m.side_id = s.id) AS number_of_memberships
	FROM sides AS s
	ORDER BY number_of_memberships DESC
	LIMIT $1`

	rows, err := ctx.Executor().QueryContext(ctx, query, limit)
	if err != nil {
		return make([]*entity.Side, 0), err
	}

	var sides []*entity.Side

	for rows.Next() {
		side := new(entity.Side)
		numberOfMemberships := 0
		if err = rows.Scan(&side.ID, &side.Nick, &side.Name, &side.Description, &side.CreatedAt, &numberOfMemberships); err != nil {
			p.logger.Errorf("%+v", err)
			continue
		}
		go p.logger.Debugf("Retrieved %+v with %d members", side, numberOfMemberships)
		sides = append(sides, side)
	}

	if err = rows.Err(); err != nil {
		return sides, err
	}

	return sides, nil
}

func (p *PgSideRepository) FindManyWithOffsetAndLimit(ctx db.DBContext[db.SqlExecutor], offset int, limit int) (entity.Sides, error) {
	var sides entity.Sides
	query := `SELECT id, nick, name, description, created_at FROM sides LIMIT $1 OFFSET $2`
	rows, err := ctx.Executor().QueryContext(ctx, query, limit, offset)
	if err != nil {
		return sides, err
	}

	for rows.Next() {
		side := new(entity.Side)
		if err := rows.Scan(&side.ID, &side.Nick, &side.Name, &side.Description, &side.CreatedAt); err != nil {
			go p.logger.Errorf("%+v", err)
			continue
		}
		sides = append(sides, side)
	}

	if err = rows.Err(); err != nil {
		return sides, err
	}

	return sides, nil
}

func (p *PgSideRepository) FindOffsetLimitedWithLargestMemberships(ctx db.DBContext[db.SqlExecutor], offset int, limit int) (entity.Sides, error) {
	query := `SELECT s.id, s.nick, s.name, s.description, s.created_at, 
    (SELECT COUNT(*) FROM memberships m WHERE m.side_id = s.id) AS number_of_memberships
	FROM sides AS s
	ORDER BY number_of_memberships DESC
	LIMIT $1
	OFFSET $2`

	rows, err := ctx.Executor().QueryContext(ctx, query, limit, offset)
	if err != nil {
		return make([]*entity.Side, 0), err
	}

	var sides []*entity.Side

	for rows.Next() {
		side := new(entity.Side)
		numberOfMemberships := 0
		if err = rows.Scan(&side.ID, &side.Nick, &side.Name, &side.Description, &side.CreatedAt, &numberOfMemberships); err != nil {
			p.logger.Errorf("%+v", err)
			continue
		}
		p.logger.Debugf("Retrieved %+v with %d members", side, numberOfMemberships)
		sides = append(sides, side)
	}

	if err = rows.Err(); err != nil {
		return sides, err
	}

	return sides, nil
}
