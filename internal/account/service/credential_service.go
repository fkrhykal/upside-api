package service

import (
	"context"
	"time"

	"github.com/fkrhykal/upside-api/internal/account/dto"
)

type CredentialService interface {
	GenerateToken(ctx context.Context, userCredential *dto.UserCredential, expiredAt time.Time) (dto.CredentialToken, error)
	RetrieveUserCredential(context.Context, dto.CredentialToken) (*dto.UserCredential, error)
}
