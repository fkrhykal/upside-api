package validation

import (
	"reflect"

	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type GoPlaygroundValidator struct {
	logger log.Logger
	v      *validator.Validate
	trans  ut.Translator
}

func (g *GoPlaygroundValidator) Validate(data any) error {
	err := g.v.Struct(data)
	if err == nil {
		return nil
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return err
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	detail := make(ErrorDetail)

	for _, err := range validationErrors {
		detail[err.Field()] = err.Translate(g.trans)
	}

	return &ValidationError{
		Detail: detail,
	}
}

func NewGoPlaygroundValidator(logger log.Logger) Validator {
	v := validator.New()

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("name")
	})

	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator(en.Locale())

	en_translations.RegisterDefaultTranslations(v, trans)

	return &GoPlaygroundValidator{
		logger: logger,
		v:      v,
		trans:  trans,
	}
}
