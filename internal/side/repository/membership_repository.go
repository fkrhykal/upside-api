package repository

import (
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/google/uuid"
)

type MembershipRepository[T any] interface {
	Save(ctx db.DBContext[T], membership *entity.Membership) error
	FindManyByMemberID(ctx db.DBContext[T], memberID uuid.UUID) ([]*entity.Membership, error)
	FindManyBySideIDsAndMemberID(ctx db.DBContext[T], sideIDs uuid.UUIDs, memberID uuid.UUID) ([]*entity.Membership, error)

	FindOffsetLimitedByMemberID(ctx db.DBContext[T], memberID uuid.UUID, offset int, limit int) (entity.Memberships, error)
}
