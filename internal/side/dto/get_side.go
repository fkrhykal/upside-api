package dto

import (
	"github.com/fkrhykal/upside-api/internal/shared/pagination"
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

type GetSidesResponse struct {
	Sides    Sides                           `json:"sides"`
	Metadata *pagination.OffsetBasedMetadata `json:"metadata"`
}
