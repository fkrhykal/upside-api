package config

import "github.com/fkrhykal/upside-api/internal/account/service"

func DefaultJwtCredentialConfig() *service.JwtCredentialConfig {
	return &service.JwtCredentialConfig{
		SignedKey: []byte(MustGetEnv("JWT_KEY")),
	}
}
