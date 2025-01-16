package app

import (
	"testing"

	"github.com/fkrhykal/upside-api/internal/shared/log"
	v "github.com/fkrhykal/upside-api/internal/shared/validation"
)

type Data struct {
	Amount  int    `validate:"min=12" name:"amount"`
	Message string `validate:"required,len=10" name:"message"`
}

func TestGoPlaygroundValidator(t *testing.T) {
	testLogger := log.NewTestLogger(t)
	validator := NewGoPlaygroundValidator(testLogger)

	data := &Data{
		Amount:  1,
		Message: "",
	}

	err := validator.Validate(data)

	validationError, ok := err.(*v.ValidationError)
	if !ok {
		t.Fatal(err)
	}

	if _, ok := validationError.Detail["amount"]; !ok {
		t.Fatal("no amount field on error")
	}

	if _, ok := validationError.Detail["message"]; !ok {
		t.Fatal("no message field on error")
	}
}
