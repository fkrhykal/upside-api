package router

import (
	"github.com/fkrhykal/upside-api/internal/account/handler"
	"github.com/fkrhykal/upside-api/internal/account/service"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/gofiber/fiber/v2"
)

func AuthRouter(logger log.Logger, authService service.AuthService) func(fiber.Router) {
	return func(app fiber.Router) {
		app.Route("/auth", func(router fiber.Router) {
			router.Post("/_sign-up", handler.SignUpHandler(logger, authService))
			router.Post("/_sign-in", handler.SignInHandler(logger, authService))
		})
	}
}
