package auth

import (
	"context"
	"time"
)

type AuthProvider interface {
	RetrieveCredential(ctx context.Context, token Token) (*UserCredential, error)
	GenerateToken(ctx context.Context, credential *UserCredential, expiredAt time.Time) (Token, error)
}
