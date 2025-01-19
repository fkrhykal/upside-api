package dto

import "github.com/google/uuid"

type UserCredential struct {
	ID uuid.UUID `json:"id"`
}

type CredentialToken string
