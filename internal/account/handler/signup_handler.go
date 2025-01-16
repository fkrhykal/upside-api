package handler

import (
	"github.com/fkrhykal/upside-api/internal/account/dto"
	"github.com/fkrhykal/upside-api/internal/account/service"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/gofiber/fiber/v2"
)

type SuccessWithData struct {
	Code    int    `json:"code" example:"201"`
	Message string `json:"message" example:"created"`
	Data    any    `json:"data"`
}

// @Summary		Sign up
// @Tags		Accounts
// @Accept		json
// @Produce		json
// @Param		SignUpRequest body dto.SignUpRequest true "Request body for sign up"
// @Router		/auth/_sign-up [post]
// @Success		201 {object} SuccessWithData{data=dto.SignUpResponse}
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
		return c.Status(fiber.StatusCreated).JSON(&SuccessWithData{
			Code:    fiber.StatusCreated,
			Message: "created",
			Data:    res,
		})
	}
}
