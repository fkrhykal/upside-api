package repository

import (
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/google/uuid"
)

type SideRepository[T any] interface {
	Save(ctx db.DBContext[T], side *entity.Side) error
	FindById(ctx db.DBContext[T], id uuid.UUID) (*entity.Side, error)
}
