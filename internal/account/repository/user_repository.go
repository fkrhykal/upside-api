package repository

import (
	"github.com/fkrhykal/upside-api/internal/account/entity"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/google/uuid"
)

type UserRepository[T any] interface {
	Save(ctx db.DBContext[T], user *entity.User) error
	FindByUsername(ctx db.DBContext[T], username string) (*entity.User, error)
	FindById(ctx db.DBContext[T], id uuid.UUID) (*entity.User, error)
}
