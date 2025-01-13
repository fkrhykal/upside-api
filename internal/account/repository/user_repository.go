package repository

import (
	"github.com/fkrhykal/upside-api/internal/account/entity"
	"github.com/fkrhykal/upside-api/internal/shared/db"
)

type UserRepository[T any] interface {
	Save(ctx db.DBContext[T], user *entity.User) error
	FindByUsername(ctx db.DBContext[T], username string) (*entity.User, error)
}
