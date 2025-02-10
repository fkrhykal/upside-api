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

func CreatePostHandler(logger log.Logger, postService service.PostService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authCtx := auth.FromFiberCtx(c)
		req := new(dto.CreatePostRequest)
		logger.Info(string(c.BodyRaw()))
		if err := c.BodyParser(req); err != nil {
			logger.Errorf("%+v", err)
			return err
		}
		id := c.Params("sideID")
		sideID, err := uuid.Parse(id)
		if err != nil {
			return response.FailureFromFiber(c, fiber.ErrNotFound)
		}
		req.SideID = sideID
		res, err := postService.CreatePost(authCtx, req)
		if err != nil {
			return err
		}
		return response.SuccessWithData(c, fiber.StatusCreated, res)
	}
}
