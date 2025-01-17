package response

import "github.com/gofiber/fiber/v2"

type Failure[T any] struct {
	Code  int `json:"code"`
	Error T   `json:"error"`
} //@name Failure

func FailureWithDetail[T any, P int](c *fiber.Ctx, code P, detail T) error {
	return c.Status(int(code)).JSON(detail)
}
