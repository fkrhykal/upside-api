package repository

import (
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/side/entity"
)

type MembershipRepository[T any] interface {
	Save(ctx db.DBContext[T], membership *entity.Membership) error
}
