package utils

import (
	"encoding/json"
	"fmt"
)

type TypeError = map[string]string

func HandleUnmarshalTypeError(err *json.UnmarshalTypeError) TypeError {
	detail := make(TypeError)
	detail[err.Field] = fmt.Sprintf("%s must be %s", err.Field, err.Type)
	return detail
}
