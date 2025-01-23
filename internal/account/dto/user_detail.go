package dto

import "github.com/google/uuid"

type UserDetail struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}
