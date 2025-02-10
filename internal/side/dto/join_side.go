package dto

import "github.com/google/uuid"

type JoinSideRequest struct {
	SideID uuid.UUID
}

type JoinSideResponse struct {
	SideID       uuid.UUID
	MembershipID uuid.UUID
}
