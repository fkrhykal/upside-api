package handler

import (
	"github.com/fkrhykal/upside-api/internal/account/dto"
	"github.com/fkrhykal/upside-api/internal/account/service"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/response"
	"github.com/gofiber/fiber/v2"
)

// @Summary		Sign in
// @Tags		Accounts
// @Accept		json
// @Produce		json
// @Param		SignInRequest body dto.SignInRequest true "Request body for sign in"
// @Router		/auth/_sign-in [post]
// @Success		200 {object} response.Success[dto.SignInResponse]
// @Failure		401 {object} response.Failure[string]
// @Failure		500 {object} response.Failure[string]
func SignInHandler(logger log.Logger, authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(dto.SignInRequest)
		c.BodyParser(req)
		res, err := authService.SignIn(c.UserContext(), req)
		if err != nil {
			return err
		}
		return response.SuccessWithData(c, fiber.StatusOK, res)
	}
}
