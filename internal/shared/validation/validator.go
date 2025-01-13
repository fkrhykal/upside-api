package validation

type ValidationError struct {
}

func (e *ValidationError) Error() string {
	return ""
}

type Validator interface {
	Validate(data any) *ValidationError
}

type ValidatorImpl struct{}

func (v *ValidatorImpl) Validate(data any) *ValidationError {
	return nil
}

func NewValidator() Validator {
	return &ValidatorImpl{}
}
