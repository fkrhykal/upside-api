package router

import (
	"github.com/fkrhykal/upside-api/internal/account/handler"
	"github.com/fkrhykal/upside-api/internal/account/service"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/gofiber/fiber/v2"
)

func AuthRouter(logger log.Logger, authService service.AuthService) func(*fiber.App) {
	return func(app *fiber.App) {
		app.Route("/auth", func(router fiber.Router) {
			router.Post("/_sign-up", handler.SignUpHandler(logger, authService))
		})
	}
}
