package router

import (
	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/side/handler"
	"github.com/fkrhykal/upside-api/internal/side/service"
	"github.com/gofiber/fiber/v2"
)

func SideRouter(logger log.Logger, sideService service.SideService, authProvider auth.AuthProvider) func(*fiber.App) {
	return func(app *fiber.App) {
		router := app.Route("/sides", func(router fiber.Router) {})

		router.Use(auth.CredentialParserMiddleware(authProvider))

		router.Get("/", handler.GetSidesHandler(logger, sideService))
		router.Post("/", handler.CreateSideHandler(logger, sideService))
	}
}
