package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Credential struct{}

var CredentialKey Credential

type FiberCredentialRegistry struct {
	*fiber.Ctx
}

func (c *FiberCredentialRegistry) SetCredential(credential *UserCredential) {
	c.Locals(CredentialKey, credential)
}

func (c *FiberCredentialRegistry) GetCredential() *UserCredential {
	return c.Locals(CredentialKey).(*UserCredential)
}

func NewFiberCredentialRegistry(c *fiber.Ctx) *FiberCredentialRegistry {
	return &FiberCredentialRegistry{Ctx: c}
}

func CredentialParserMiddleware(authProvider AuthProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		registry := &FiberCredentialRegistry{Ctx: c}
		authorizationHeader := c.Get("Authorization")
		token, ok := strings.CutPrefix(authorizationHeader, "Bearer ")
		if !ok {
			return registry.Next()
		}
		credential, err := authProvider.RetrieveCredential(c.UserContext(), Token(token))
		if err != nil {
			return registry.Next()
		}
		registry.SetCredential(credential)
		return registry.Next()
	}
}
