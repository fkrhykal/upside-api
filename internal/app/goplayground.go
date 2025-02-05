package app

import (
	"reflect"
	"unicode"

	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/validation"
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

func NewGoPlaygroundValidator(logger log.Logger) validation.Validator {
	v := validator.New()
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator(en.Locale())
	en_translations.RegisterDefaultTranslations(v, trans)

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		json := field.Tag.Get("json")
		if json == "-" {
			return ""
		}
		return json
	})

	v.RegisterValidation("password", PasswordValidation(logger))

	v.RegisterTranslation("password", trans, func(ut ut.Translator) error {
		return ut.Add(
			"password",
			"{0} must contain uppercase letter, lowercase letter, number, special character and no space",
			true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password", fe.Field())
		return t
	})

	return &GoPlaygroundValidator{
		logger: logger,
		v:      v,
		trans:  trans,
	}
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

	detail := make(validation.ErrorDetail)

	for _, err := range validationErrors {
		detail[err.Field()] = err.Translate(g.trans)
	}

	return &validation.ValidationError{
		Detail: detail,
	}
}

func PasswordValidation(logger log.Logger) validator.Func {
	return func(fl validator.FieldLevel) bool {
		password := fl.Field().String()

		logger.Debugf("validating password %s", password)

		hasDigit := false
		hasLower := false
		hasUpper := false
		hasSymbol := false

		for _, c := range password {
			if !hasLower && unicode.IsLower(c) {
				hasLower = true
				continue
			}
			if !hasUpper && unicode.IsUpper(c) {
				hasUpper = true
				continue
			}
			if !hasDigit && unicode.IsDigit(c) {
				hasDigit = true
				continue
			}
			if !hasSymbol && validSymbols.Exist(c) {
				hasSymbol = true
				continue
			}
			if unicode.IsSpace(c) {
				logger.Debugf("password has space")
				return false
			}
		}

		logger.Debugf("password has lowercase: %+v", hasLower)
		logger.Debugf("password has uppercase: %+v", hasUpper)
		logger.Debugf("password has symbol: %+v", hasSymbol)
		logger.Debugf("password has number: %+v", hasDigit)

		result := hasLower && hasUpper && hasDigit && hasSymbol

		logger.Debugf("password validation result: %+v", result)

		return result
	}
}

type ValidSymbolRegistry map[string]struct{}

func (s ValidSymbolRegistry) Exist(key rune) bool {
	_, ok := s[string(key)]
	return ok
}

var validSymbols = ValidSymbolRegistry{
	"!":  {},
	"@":  {},
	"#":  {},
	"$":  {},
	"%":  {},
	"^":  {},
	"&":  {},
	"*":  {},
	"(":  {},
	")":  {},
	"-":  {},
	"_":  {},
	"=":  {},
	"+":  {},
	"[":  {},
	"]":  {},
	"{":  {},
	"}":  {},
	"\\": {},
	"|":  {},
	";":  {},
	":":  {},
	"'":  {},
	"\"": {},
	",":  {},
	".":  {},
	"/":  {},
	"?":  {},
	"<":  {},
	">":  {},
	"~":  {},
	"`":  {},
}
