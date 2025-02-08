package auth

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type AuthContext struct {
	context.Context
	Credential *UserCredential
}

func (a *AuthContext) Authenticated() bool {
	return a.Credential != nil
}

func NewAuthContext(ctx context.Context, credential *UserCredential) *AuthContext {
	return &AuthContext{
		Context:    ctx,
		Credential: credential,
	}
}

func FromFiberCtx(c *fiber.Ctx) *AuthContext {
	registry := NewFiberCredentialRegistry(c)
	return &AuthContext{
		Context:    c.UserContext(),
		Credential: registry.GetCredential(),
	}
}
