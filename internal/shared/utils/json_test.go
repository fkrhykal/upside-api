package utils_test

import (
	"encoding/json"
	"testing"

	"github.com/fkrhykal/upside-api/internal/shared/utils"
)

type Data struct {
	Message string   `json:"message"`
	Items   []string `json:"items"`
	Amount  int      `json:"amount"`
}

func TestHandleDataMessageUnmarshalTypeError(t *testing.T) {

	data := new(Data)
	payload := `{"message": 4, "items": ""}`

	err := json.Unmarshal([]byte(payload), data)

	if err, ok := err.(*json.UnmarshalTypeError); !ok {
		t.Fatal(err)
	}

	detail := utils.HandleUnmarshalTypeError(err.(*json.UnmarshalTypeError))

	message, ok := detail["message"]
	if !ok {
		t.Fatal("Error with key message should be exist")
	}
	if message != "message must be string, but found number" {
		t.Fatalf("Error message mismatch: %s", message)
	}
}

func TestHandleDataItemsUnmarshalTypeError(t *testing.T) {

	data := new(Data)
	payload := `{"message": "d", "items": ""}`

	err := json.Unmarshal([]byte(payload), data)

	if err, ok := err.(*json.UnmarshalTypeError); !ok {
		t.Fatal(err)
	}

	detail := utils.HandleUnmarshalTypeError(err.(*json.UnmarshalTypeError))

	message, ok := detail["items"]
	if !ok {
		t.Fatal("Error with key message should be exist")
	}
	if message != "items must be []string, but found string" {
		t.Fatalf("Error message mismatch: %s", message)
	}
}

func TestHandleDataAmountUnmarshalTypeError(t *testing.T) {

	data := new(Data)
	payload := `{"message": "d", "items": [], "amount": []}`

	err := json.Unmarshal([]byte(payload), data)

	if err, ok := err.(*json.UnmarshalTypeError); !ok {
		t.Fatal(err)
	}

	detail := utils.HandleUnmarshalTypeError(err.(*json.UnmarshalTypeError))

	message, ok := detail["amount"]
	if !ok {
		t.Fatal("Error with key message should be exist")
	}
	if message != "amount must be integer, but found array" {
		t.Fatalf("Error message mismatch: %s", message)
	}
}
