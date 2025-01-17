package response

import "github.com/gofiber/fiber/v2"

type Success[T any] struct {
	Code int `json:"code"`
	Data T   `json:"data"`
} //@name Success

func SuccessWithData[T any, P int](c *fiber.Ctx, code P, data T) error {
	return c.Status(int(code)).JSON(data)
}
