package repository

import (
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/pagination"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type PostRepository[T any] interface {
	Save(ctx db.DBContext[T], post *entity.Post) error
	FindByID(ctx db.DBContext[T], postID ulid.ULID) (*entity.Post, error)
	FindManyWithULIDCursor(ctx db.DBContext[T], cursor pagination.ULIDCursor) (*pagination.ULIDCursorMetadata[*entity.Post], error)
	FindManyWithULIDCursorInSides(ctx db.DBContext[T], cursor pagination.ULIDCursor, sideIDs uuid.UUIDs) (*pagination.ULIDCursorMetadata[*entity.Post], error)
}
