package utils_test

import (
	"encoding/json"
	"testing"

	"github.com/fkrhykal/upside-api/internal/shared/utils"
)

type Data struct {
	Message string `json:"message"`
}

func TestHandleUnmarshalTypeError(t *testing.T) {

	data := new(Data)
	payload := `{"message": 4}`

	err := json.Unmarshal([]byte(payload), data)

	if err, ok := err.(*json.UnmarshalTypeError); !ok {
		t.Fatal(err)
	}

	detail := utils.HandleUnmarshalTypeError(err.(*json.UnmarshalTypeError))

	message, ok := detail["message"]
	if !ok {
		t.Fatal("Error with key message should be exist")
	}
	if message != "message must be string" {
		t.Fatalf("Error message mismatch: %s", message)
	}
}
