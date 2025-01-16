package dto_test

import (
	"testing"

	"github.com/fkrhykal/upside-api/internal/account/dto"
	"github.com/fkrhykal/upside-api/internal/app"
	"github.com/fkrhykal/upside-api/internal/shared/validation"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/suite"
)

type SignUpRequestDtoValidationSuite struct {
	app.GoPlaygroundValidationSuite
}

func (s *SignUpRequestDtoValidationSuite) TestRequestValid() {
	req := &dto.SignUpRequest{
		Username: faker.Username(),
		Password: "fbefei&%24HH+",
	}
	s.T().Log(req.Password)
	if err := s.Validator.Validate(req); err != nil {
		s.T().Fatal(err)
	}
}

func (s *SignUpRequestDtoValidationSuite) TestUsernameEmpty() {
	req := dto.SignUpRequest{
		Password: faker.Password(),
	}
	err := s.Validator.Validate(req)
	if err == nil {
		s.T().Fatal()
	}
	validationError, ok := err.(*validation.ValidationError)
	if !ok {
		s.T().Fatal(err)
	}
	if _, ok := validationError.Detail["username"]; !ok {
		s.T().Fatal()
	}
}

func (s *SignUpRequestDtoValidationSuite) TestUsernameInvalid() {
	req := dto.SignUpRequest{
		Username: "sdfdsnf*&#",
		Password: faker.Password(),
	}
	err := s.Validator.Validate(req)
	if err == nil {
		s.T().Fatal()
	}
	validationError, ok := err.(*validation.ValidationError)
	if !ok {
		s.T().Fatal(err)
	}
	if _, ok := validationError.Detail["username"]; !ok {
		s.T().Fatal()
	}
}

func (s *SignUpRequestDtoValidationSuite) TestPasswordEmpty() {
	req := dto.SignUpRequest{
		Username: faker.Username(),
	}
	err := s.Validator.Validate(req)
	if err == nil {
		s.T().Fatal()
	}
	validationError, ok := err.(*validation.ValidationError)
	if !ok {
		s.T().Fatal(err)
	}
	if _, ok := validationError.Detail["password"]; !ok {
		s.T().Fatal()
	}
}

func (s *SignUpRequestDtoValidationSuite) TestPasswordInvalid() {
	req := dto.SignUpRequest{
		Username: faker.Username(),
		Password: "sdndfnksdf",
	}
	err := s.Validator.Validate(req)
	if err == nil {
		s.T().Fatal()
	}
	validationError, ok := err.(*validation.ValidationError)
	if !ok {
		s.T().Fatal(err)
	}
	if _, ok := validationError.Detail["password"]; !ok {
		s.T().Fatal()
	}
}

func TestSignUpRequestValidation(t *testing.T) {
	suite.Run(t, new(SignUpRequestDtoValidationSuite))
}
