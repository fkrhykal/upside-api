package service

import (
	"context"
	"time"

	"github.com/fkrhykal/upside-api/internal/account/dto"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/golang-jwt/jwt/v5"
)

type JwtCredentialConfig struct {
	SignedKey []byte
}

type AuthenticationClams struct {
	jwt.RegisteredClaims
	UserCredential *dto.UserCredential `json:"userCredential"`
}

type JwtCredentialService struct {
	logger log.Logger
	config *JwtCredentialConfig
}

func (s *JwtCredentialService) GenerateToken(ctx context.Context, credential *dto.UserCredential, expiredAt time.Time) (dto.CredentialToken, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, *&AuthenticationClams{
		UserCredential: credential,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "upside",
			ExpiresAt: jwt.NewNumericDate(expiredAt),
		},
	})
	signedJwt, err := jwtToken.SignedString(s.config.SignedKey)
	return dto.CredentialToken(signedJwt), err
}

func (s *JwtCredentialService) RetrieveUserCredential(ctx context.Context, token dto.CredentialToken) (*dto.UserCredential, error) {
	jwtToken, err := jwt.ParseWithClaims(string(token), &AuthenticationClams{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "upside",
		},
	}, func(t *jwt.Token) (interface{}, error) {
		return s.config.SignedKey, nil
	})
	if err != nil {
		return nil, err
	}
	return jwtToken.Claims.(*AuthenticationClams).UserCredential, nil
}

func NewJwtCredentialService(logger log.Logger, config *JwtCredentialConfig) *JwtCredentialService {
	return &JwtCredentialService{
		logger: logger,
		config: config,
	}
}
