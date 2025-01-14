package validation

import (
	"fmt"
)

type ErrorDetail = map[string]string

type ValidationError struct {
	Detail ErrorDetail
}

func (e *ValidationError) Error() string {
	return fmt.Sprint(e.Detail)
}

type Validator interface {
	Validate(data any) error
}
