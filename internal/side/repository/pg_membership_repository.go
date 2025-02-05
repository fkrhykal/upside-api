package repository

import (
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/side/entity"
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
