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

type PgMembershipRepository struct {
	logger log.Logger
}

func NewPgMembershipRepository(logger log.Logger) MembershipRepository[db.SqlExecutor] {
	return &PgMembershipRepository{
		logger: logger,
	}
}

func (m *PgMembershipRepository) Save(ctx db.DBContext[db.SqlExecutor], membership *entity.Membership) error {
	query := "INSERT INTO memberships(id, member_id, side_id, role) VALUES($1, $2, $3, $4)"
	_, err := ctx.Executor().ExecContext(ctx, query, membership.ID, membership.Member, membership.Side, membership.Role)
	if err != nil {
		return err
	}
	return nil
}

func (m *PgMembershipRepository) FindManyByMemberID(ctx db.DBContext[db.SqlExecutor], memberID uuid.UUID) ([]*entity.Membership, error) {

	query := "SELECT id, member_id, side_id, role FROM memberships WHERE member_id = $1"

	rows, err := ctx.Executor().QueryContext(ctx, query, memberID)

	var memberships []*entity.Membership

	if err != nil {
		return memberships, err
	}

	for rows.Next() {
		membership := new(entity.Membership)
		if err := rows.Scan(&membership.ID, &membership.Member, &membership.Side, &membership.Role); err != nil {
			continue
		}
		memberships = append(memberships, membership)
	}

	if err = rows.Err(); err != nil {
		return memberships, err
	}

	return memberships, nil
}

func (m *PgMembershipRepository) FindManyBySideIDsAndMemberID(ctx db.DBContext[db.SqlExecutor], sideIDs uuid.UUIDs, memberID uuid.UUID) ([]*entity.Membership, error) {
	query := `SELECT id, member_id, side_id, role FROM memberships WHERE side_id = ANY($1) AND member_id = $2`
	rows, err := ctx.Executor().QueryContext(ctx, query, pq.Array(sideIDs), memberID)
	var memberships []*entity.Membership
	if err != nil {
		return memberships, nil
	}
	for rows.Next() {
		membership := new(entity.Membership)
		if err := rows.Scan(&membership.ID, &membership.Member, &membership.Side, &membership.Role); err != nil {
			continue
		}
		memberships = append(memberships, membership)
	}
	if err = rows.Err(); err != nil {
		return memberships, err
	}
	return memberships, nil
}

func (m *PgMembershipRepository) FindOffsetLimitedByMemberID(ctx db.DBContext[db.SqlExecutor], memberID uuid.UUID, offset int, limit int) (entity.Memberships, error) {
	query := `SELECT id, member_id, side_id, role FROM memberships WHERE member_id = $1 OFFSET $2 LIMIT $3`
	rows, err := ctx.Executor().QueryContext(ctx, query, memberID, offset, limit)
	var memberships []*entity.Membership
	if err != nil {
		return memberships, err
	}
	for rows.Next() {
		membership := new(entity.Membership)
		if err := rows.Scan(&membership.ID, &membership.Member, &membership.Side, &membership.Role); err != nil {
			continue
		}
		memberships = append(memberships, membership)
	}
	if err = rows.Err(); err != nil {
		return memberships, err
	}
	return memberships, nil
}

func (m *PgMembershipRepository) FindBySideIDAndMemberID(ctx db.DBContext[db.SqlExecutor], memberID uuid.UUID, sideID uuid.UUID) (*entity.Membership, error) {
	query := `SELECT id, member_id, side_id, role FROM memberships WHERE member_id = $1 AND side_id = $2`
	membership := new((entity.Membership))
	err := ctx.Executor().
		QueryRowContext(ctx, query, memberID, sideID).
		Scan(&membership.ID, &membership.Member, &membership.Side, &membership.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return membership, nil
}
