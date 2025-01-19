package response

import "github.com/gofiber/fiber/v2"

type Failure[T any] struct {
	Code  int `json:"code"`
	Error T   `json:"error"`
} //@name Failure

func FailureWithDetail[T any](c *fiber.Ctx, code int, detail T) error {
	return c.Status(int(code)).JSON(&Failure[T]{
		Code:  code,
		Error: detail,
	})
}

func FailureFromFiber(c *fiber.Ctx, err *fiber.Error) error {
	return c.Status(err.Code).JSON(&Failure[string]{
		Code:  err.Code,
		Error: err.Message,
	})
}
