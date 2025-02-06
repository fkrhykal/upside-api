package auth

import (
	"strings"

	"github.com/fkrhykal/upside-api/internal/shared/response"
	"github.com/gofiber/fiber/v2"
)

type Credential struct{}

var CredentialKey Credential

type AuthCtx struct {
	*fiber.Ctx
}

func (c *AuthCtx) SetCredential(credential *UserCredential) {
	c.Locals(CredentialKey, credential)
}

func (c *AuthCtx) GetCredential() *UserCredential {
	return c.Locals(CredentialKey).(*UserCredential)
}

func FromCtx(c *fiber.Ctx) *AuthCtx {
	return &AuthCtx{Ctx: c}
}

func AuthMiddleware(authProvider AuthProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authCtx := &AuthCtx{Ctx: c}
		authorizationHeader := authCtx.Get("Authorization")
		token, ok := strings.CutPrefix(authorizationHeader, "Bearer ")
		if !ok {
			return response.FailureFromFiber(c, fiber.ErrUnauthorized)
		}
		credential, err := authProvider.RetrieveCredential(c.UserContext(), Token(token))
		if err != nil {
			return response.FailureFromFiber(c, fiber.ErrUnauthorized)
		}
		authCtx.SetCredential(credential)
		return authCtx.Next()
	}
}
