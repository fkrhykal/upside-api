package entity

import (
	"github.com/google/uuid"
)

type Side struct {
	ID          uuid.UUID
	Nick        string
	Name        string
	Description string
	CreatedAt   uint64
}
