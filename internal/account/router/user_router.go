package router

import (
	"github.com/fkrhykal/upside-api/internal/account/handler"
	"github.com/fkrhykal/upside-api/internal/account/service"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/gofiber/fiber/v2"
)

func UserRouter(logger log.Logger, userService service.UserService) func(fiber.Router) {
	return func(app fiber.Router) {
		app.Route("/users", func(router fiber.Router) {
			router.Get("/:id", handler.GetUserDetailHandler(logger, userService))
		})
	}
}
