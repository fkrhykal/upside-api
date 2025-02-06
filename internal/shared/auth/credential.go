package auth

import "github.com/google/uuid"

type UserCredential struct {
	ID uuid.UUID
}

type Token string
