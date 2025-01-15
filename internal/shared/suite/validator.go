package suite

import (
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/validation"
	"github.com/stretchr/testify/suite"
)

type GoPlaygroundValidationSuite struct {
	suite.Suite
	Validator validation.Validator
}

func (s *GoPlaygroundValidationSuite) SetupSuite() {
	logger := log.NewTestLogger(s.T())
	s.Validator = validation.NewGoPlaygroundValidator(logger)
}
