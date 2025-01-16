package app

import (
	"encoding/json"

	"github.com/fkrhykal/upside-api/internal/shared/log"
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

type DetailedError struct {
	*fiber.Error
	Detail any `json:"error"`
}

func ErrorWithDetail(c *fiber.Ctx, err *fiber.Error, detail any) error {
	detailedError := &DetailedError{
		Error:  err,
		Detail: detail,
	}
	return c.Status(err.Code).JSON(detailedError)
}

func SetupErrorHandler(logger log.Logger) fiber.ErrorHandler {

	return func(c *fiber.Ctx, err error) error {
		switch err.(type) {
		case *validation.ValidationError:
			return ErrorWithDetail(c, fiber.ErrBadRequest, err.(*validation.ValidationError).Detail)
		case *fiber.Error:
			return c.SendStatus(err.(*fiber.Error).Code)
		case *json.UnmarshalTypeError:
			detail := utils.HandleUnmarshalTypeError(err.(*json.UnmarshalTypeError))
			return ErrorWithDetail(c, fiber.ErrBadRequest, detail)
		default:
			logger.Errorf("Unpredicted error: %+v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}
}
