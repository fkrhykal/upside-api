package entity

import "github.com/google/uuid"

type Role string

const (
	FOUNDER Role = "founder"
	ADMIN   Role = "admin"
	MEMBER  Role = "member"
)

type Membership struct {
	ID     uuid.UUID
	Member uuid.UUID
	Side   uuid.UUID
	Role   Role
}
