package validation

import (
	"fmt"
)

type ErrorDetail = map[string]string

type ValidationError struct {
	Detail ErrorDetail
}

func (e *ValidationError) Exist(key string) bool {
	_, ok := e.Detail[key]
	return ok
}

func (e *ValidationError) Add(key string, message string) {
	e.Detail[key] = message
}

func (e *ValidationError) Error() string {
	return fmt.Sprint(e.Detail)
}

type Validator interface {
	Validate(data any) error
}
