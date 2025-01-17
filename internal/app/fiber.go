package app

import (
	"encoding/json"

	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/response"
	"github.com/fkrhykal/upside-api/internal/shared/utils"
	"github.com/fkrhykal/upside-api/internal/shared/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func NewFiber(logger log.Logger) *fiber.App {

	app := fiber.New(fiber.Config{
		ErrorHandler: SetupErrorHandler(logger),
	})

	app.Use(cors.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}

func SetupErrorHandler(logger log.Logger) fiber.ErrorHandler {

	return func(c *fiber.Ctx, err error) error {
		switch err.(type) {
		case *validation.ValidationError:
			detail := err.(*validation.ValidationError).Detail
			return response.FailureWithDetail(c, fiber.StatusBadRequest, detail)
		case *fiber.Error:
			return c.SendStatus(err.(*fiber.Error).Code)
		case *json.UnmarshalTypeError:
			detail := utils.HandleUnmarshalTypeError(err.(*json.UnmarshalTypeError))
			return response.FailureWithDetail(c, fiber.StatusBadRequest, detail)
		default:
			logger.Errorf("Unpredicted error: %+v", err)
			return response.FailureWithDetail(c, fiber.StatusBadRequest, fiber.ErrBadRequest.Message)
		}
	}
}
