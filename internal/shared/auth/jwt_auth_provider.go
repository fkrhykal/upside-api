package auth

import (
	"context"
	"time"

	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/golang-jwt/jwt/v5"
)

type JwtAuthConfig struct {
	SignedKey []byte
}

type AuthenticationClams struct {
	jwt.RegisteredClaims
	UserCredential *UserCredential `json:"userCredential"`
}

type JwtAuthProvider struct {
	logger log.Logger
	config *JwtAuthConfig
}

func (s *JwtAuthProvider) GenerateToken(ctx context.Context, credential *UserCredential, expiredAt time.Time) (Token, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, *&AuthenticationClams{
		UserCredential: credential,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "upside",
			ExpiresAt: jwt.NewNumericDate(expiredAt),
		},
	})
	signedJwt, err := jwtToken.SignedString(s.config.SignedKey)
	return Token(signedJwt), err
}

func (s *JwtAuthProvider) RetrieveCredential(ctx context.Context, token Token) (*UserCredential, error) {
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

func NewJwtAuthProvider(logger log.Logger, config *JwtAuthConfig) AuthProvider {
	return &JwtAuthProvider{
		logger: logger,
		config: config,
	}
}
