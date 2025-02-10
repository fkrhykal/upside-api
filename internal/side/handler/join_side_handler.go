package handler

import (
	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/response"
	"github.com/fkrhykal/upside-api/internal/side/dto"
	"github.com/fkrhykal/upside-api/internal/side/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @Summary		Join Side
// @Tags		Sides
// @Produce		json
// @Router		/sides/{sideID}/_join [post]
// @Param		sideID path string true "side id"
// @Security	BearerAuth
// @Success		200 {object} response.Success[dto.JoinSideResponse]
// @Failure		401 {object} response.Failure[string]
// @Failure		409 {object} response.Failure[string]
// @Failure		500 {object} response.Failure[string]
func JoinSideHandler(logger log.Logger, sideService service.SideService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("sideID")
		sideID, err := uuid.Parse(id)
		if err != nil {
			return response.FailureFromFiber(c, fiber.ErrNotFound)
		}
		authCtx := auth.FromFiberCtx(c)
		res, err := sideService.JoinSide(authCtx, &dto.JoinSideRequest{SideID: sideID})
		if err != nil {
			return err
		}
		return response.SuccessWithData(c, fiber.StatusOK, res)
	}
}
