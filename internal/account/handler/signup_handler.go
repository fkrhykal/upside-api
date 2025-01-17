package handler

import (
	"github.com/fkrhykal/upside-api/internal/account/dto"
	"github.com/fkrhykal/upside-api/internal/account/service"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/response"
	"github.com/gofiber/fiber/v2"
)

// @Summary		Sign up
// @Tags		Accounts
// @Accept		json
// @Produce		json
// @Param		SignUpRequest body dto.SignUpRequest true "Request body for sign up"
// @Router		/auth/_sign-up [post]
// @Success		201 {object} response.Success[dto.SignUpResponse]
// @Failure		400 {object} response.Failure[validation.ErrorDetail]
// @Failure		500 {object} response.Failure[string]
func SignUpHandler(logger log.Logger, authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(dto.SignUpRequest)
		if err := c.BodyParser(req); err != nil {
			return err
		}
		res, err := authService.SignUp(c.UserContext(), req)
		if err != nil {
			return err
		}
		return response.SuccessWithData(c, fiber.StatusCreated, res)
	}
}
