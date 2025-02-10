package repository

import (
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/side/entity"
)

type PostRepository[T any] interface {
	Save(ctx db.DBContext[T], post *entity.Post) error
}
