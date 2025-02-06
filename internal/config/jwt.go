package config

import (
	"github.com/fkrhykal/upside-api/internal/shared/auth"
)

func DefaultJwtAuthConfig() *auth.JwtAuthConfig {
	return &auth.JwtAuthConfig{
		SignedKey: []byte(MustGetEnv("JWT_KEY")),
	}
}
