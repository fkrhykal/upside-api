package handler

import (
	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/response"
	"github.com/fkrhykal/upside-api/internal/side/dto"
	"github.com/fkrhykal/upside-api/internal/side/service"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

func RevokeVotePostHandler(logger log.Logger, voteService service.VoteService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authCtx := auth.FromFiberCtx(c)
		postID := c.Params("postID")
		postULID, err := ulid.Parse(postID)
		if err != nil {
			return response.FailureFromFiber(c, fiber.ErrNotFound)
		}
		err = voteService.RevokeVote(authCtx, &dto.RevokeVoteRequest{
			PostID: postULID,
		})
		if err != nil {
			return err
		}
		return response.SuccessWithData(c, fiber.StatusOK, struct{}{})
	}
}
