package router

import (
	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/side/handler"
	"github.com/fkrhykal/upside-api/internal/side/service"
	"github.com/gofiber/fiber/v2"
)

func PostRouter(logger log.Logger, authProvider auth.AuthProvider, postService service.PostService, voteService service.VoteService) func(fiber.Router) {
	return func(app fiber.Router) {
		app.Route("/posts", func(router fiber.Router) {
			router.Use(auth.JwtBearerParserMiddleware(authProvider))

			router.Get("/", handler.GetPostsHandler(logger, postService))
			voteRouter := router.Group("/:postID", auth.AuthenticationMiddleware)
			voteRouter.Post("/_upvote", handler.UpVotePostHandler(logger, voteService))
			voteRouter.Post("/_downvote", handler.DownVotePostHandler(logger, voteService))
			voteRouter.Post("/_revoke-vote", handler.RevokeVotePostHandler(logger, voteService))
		})
	}
}
