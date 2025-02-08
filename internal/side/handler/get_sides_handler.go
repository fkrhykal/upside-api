package handler

import (
	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/pagination"
	"github.com/fkrhykal/upside-api/internal/shared/response"
	"github.com/fkrhykal/upside-api/internal/side/service"
	"github.com/gofiber/fiber/v2"
)

func GetSidesHandler(logger log.Logger, sideService service.SideService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authCtx := auth.FromFiberCtx(c)
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 10)
		offsetPagination := pagination.SafeOffsetBased(page, limit)
		switch c.Query("filter") {
		case "joined":
			res, err := sideService.GetJoinedSides(authCtx, offsetPagination)
			if err != nil {
				return err
			}
			return response.SuccessWithData(c, fiber.StatusOK, res)
		case "popular":
			res, err := sideService.GetPopularSides(authCtx, offsetPagination)
			if err != nil {
				return err
			}
			return response.SuccessWithData(c, fiber.StatusOK, res)
		default:
			res, err := sideService.GetSides(authCtx, offsetPagination)
			if err != nil {
				return err
			}
			return response.SuccessWithData(c, fiber.StatusOK, res)
		}
	}
}
