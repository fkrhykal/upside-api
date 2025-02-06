package handler

import (
	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/response"
	"github.com/fkrhykal/upside-api/internal/side/dto"
	"github.com/fkrhykal/upside-api/internal/side/service"
	"github.com/gofiber/fiber/v2"
)

// @Summary		Create Side
// @Tags		Sides
// @Produce		json
// @Router		/sides [post]
// @Param		CreateSideRequest body dto.CreateSideRequest true "Request body for create side"
// @Security	BearerAuth
// @Success		201 {object} response.Success[dto.CreateSideResponse]
// @Failure		401 {object} response.Failure[string]
// @Failure		500 {object} response.Failure[string]
func CreateSideHandler(logger log.Logger, sideService service.SideService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(dto.CreateSideRequest)

		if err := c.BodyParser(req); err != nil {
			return err
		}

		logger.Debug(req)

		req.FounderID = auth.FromCtx(c).GetCredential().ID

		res, err := sideService.CreateSide(c.UserContext(), req)
		if err != nil {
			return err
		}
		return response.SuccessWithData(c, fiber.StatusCreated, res)
	}
}
