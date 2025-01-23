package handler

import (
	"github.com/fkrhykal/upside-api/internal/account/service"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @Summary		Get User Detail
// @Tags		Users
// @Produce		json
// @Param		id path string true "user id" example(28fd7c57-ffde-4b4b-83c3-4781d93c268e)
// @Router		/users/{id} [get]
// @Success		200 {object} response.Success[dto.UserDetail]
// @Failure		404 {object} response.Failure[string]
// @Failure		500 {object} response.Failure[string]
func GetUserDetailHandler(logger log.Logger, userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return response.FailureFromFiber(c, fiber.ErrNotFound)
		}
		userId, err := uuid.Parse(id)
		if err != nil {
			return response.FailureFromFiber(c, fiber.ErrNotFound)
		}
		detail, err := userService.GetUserDetail(c.UserContext(), userId)
		if err != nil {
			return err
		}
		return response.SuccessWithData(c, int(fiber.StatusOK), detail)
	}
}
