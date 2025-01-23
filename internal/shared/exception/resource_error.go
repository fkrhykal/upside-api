package exception

import "fmt"

type ResourceNotFoundError struct {
	ResourceName string
}

func (nf *ResourceNotFoundError) Error() string {
	return fmt.Sprintf("%s not found", nf.ResourceName)
}

var (
	UserNotFound = &ResourceNotFoundError{ResourceName: "user"}
)
