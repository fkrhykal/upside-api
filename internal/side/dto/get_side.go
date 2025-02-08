package dto

import (
	"github.com/google/uuid"
)

type MembershipDetail struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"`
}

type Side struct {
	ID               uuid.UUID         `json:"id"`
	Nick             string            `json:"nick"`
	Name             string            `json:"name"`
	MembershipDetail *MembershipDetail `json:"membershipDetail"`
}

type Sides []*Side

var EmptySides = make(Sides, 0)
