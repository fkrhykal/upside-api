package router

import (
	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/side/handler"
	"github.com/fkrhykal/upside-api/internal/side/service"
	"github.com/gofiber/fiber/v2"
)

func SideRouter(logger log.Logger, authProvider auth.AuthProvider, sideService service.SideService, postService service.PostService) func(fiber.Router) {
	return func(app fiber.Router) {
		router := app.Route("/sides", func(router fiber.Router) {})

		router.Use(auth.JwtBearerParserMiddleware(authProvider))

		router.Get("/", handler.GetSidesHandler(logger, sideService))
		router.Post("/", handler.CreateSideHandler(logger, sideService))
		router.Post("/:sideID/_join", auth.AuthenticationMiddleware, handler.JoinSideHandler(logger, sideService))
		router.Post("/:sideID/posts", auth.AuthenticationMiddleware, handler.CreatePostHandler(logger, postService))
	}
}
