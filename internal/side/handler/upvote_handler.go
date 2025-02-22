package handler

import (
	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/response"
	"github.com/fkrhykal/upside-api/internal/side/dto"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/fkrhykal/upside-api/internal/side/service"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

func UpVotePostHandler(logger log.Logger, voteService service.VoteService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authCtx := auth.FromFiberCtx(c)
		id := c.Params("postID")
		postID, err := ulid.Parse(id)
		if err != nil {
			return response.FailureFromFiber(c, fiber.ErrNotFound)
		}
		res, err := voteService.Vote(authCtx, &dto.VoteRequest{
			PostID:   postID,
			VoteKind: entity.UpVote,
		})
		if err != nil {
			return err
		}
		return response.SuccessWithData(c, fiber.StatusOK, res)
	}
}
