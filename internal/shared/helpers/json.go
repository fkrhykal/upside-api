package helpers

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type TypeErrorDetail = map[string]string

var typeMappings = map[string]string{
	"int":  "integer",
	"bool": "boolean",
}

func mapType(t reflect.Type) string {
	if friendlyType, exists := typeMappings[t.String()]; exists {
		return friendlyType
	}
	return t.String()
}

func HandleUnmarshalTypeError(err *json.UnmarshalTypeError) TypeErrorDetail {
	detail := make(TypeErrorDetail)

	friendlyType := mapType(err.Type)
	detail[err.Field] = fmt.Sprintf("%s must be %s, but found %s", err.Field, friendlyType, err.Value)
	return detail
}
