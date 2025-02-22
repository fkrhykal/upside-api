package exception

import "fmt"

type NotFoundError struct {
	ResourceName string
}

func (nf *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", nf.ResourceName)
}

var (
	ErrUserNotFound = &NotFoundError{ResourceName: "user"}
	ErrSideNotFound = &NotFoundError{ResourceName: "side"}
	ErrPostNotFound = &NotFoundError{ResourceName: "post"}
)
