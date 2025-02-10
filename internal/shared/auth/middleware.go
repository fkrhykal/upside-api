package auth

import (
	"strings"

	"github.com/fkrhykal/upside-api/internal/shared/response"
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
	credential, ok := c.Locals(CredentialKey).(*UserCredential)
	if !ok {
		return nil
	}
	return credential
}

func NewFiberCredentialRegistry(c *fiber.Ctx) *FiberCredentialRegistry {
	return &FiberCredentialRegistry{Ctx: c}
}

func JwtBearerParserMiddleware(authProvider AuthProvider) fiber.Handler {
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

// Check for credential, return 401 if not exist
func AuthenticationMiddleware(c *fiber.Ctx) error {
	authCtx := FromFiberCtx(c)
	if !authCtx.Authenticated() {
		return response.FailureFromFiber(c, fiber.ErrUnauthorized)
	}
	return c.Next()
}
