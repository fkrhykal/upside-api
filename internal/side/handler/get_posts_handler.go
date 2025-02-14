package handler

import (
	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/pagination"
	"github.com/fkrhykal/upside-api/internal/shared/response"
	"github.com/fkrhykal/upside-api/internal/side/service"
	"github.com/gofiber/fiber/v2"
)

func GetPostsHandler(logger log.Logger, postService service.PostService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authCtx := auth.FromFiberCtx(c)
		cursor := c.Query("cursor")
		limit := c.QueryInt("limit")
		filter := c.Query("filter")
		ulidCursor, err := pagination.ParseULIDCursor(cursor, limit)
		if err != nil {
			return err
		}
		if filter == "subscribed" {
			res, err := postService.GetSubscribedPosts(authCtx, ulidCursor)
			if err != nil {
				return err
			}
			return response.SuccessWithData(c, fiber.StatusOK, res)
		}
		res, err := postService.GetLatestPosts(authCtx, ulidCursor)
		if err != nil {
			return err
		}
		return response.SuccessWithData(c, fiber.StatusOK, res)
	}
}
