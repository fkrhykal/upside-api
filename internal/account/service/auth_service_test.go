package service_test

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/suite"

	"github.com/fkrhykal/upside-api/internal/account/dto"
	"github.com/fkrhykal/upside-api/internal/account/repository"
	"github.com/fkrhykal/upside-api/internal/account/service"
	"github.com/fkrhykal/upside-api/internal/account/utils"
	"github.com/fkrhykal/upside-api/internal/app"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/log"
)

type AuthServiceTestSuite struct {
	app.PostgresContainerSuite
	authService    service.AuthService
	userRepository repository.UserRepository[db.SqlExecutor]
	ctxManager     db.CtxManager[db.SqlExecutor]
}

func (s *AuthServiceTestSuite) TestSignUp() {

	req := &dto.SignUpRequest{
		Username: faker.Username(),
		Password: "TestPassword123^$&",
	}

	ctx := context.Background()

	res, err := s.authService.SignUp(ctx, req)
	if err != nil {
		s.FailNow("Failed to sign up: ", err)
	}

	dbCtx := s.ctxManager.NewDBContext(ctx)
	user, err := s.userRepository.FindByUsername(dbCtx, req.Username)
	if err != nil {
		s.FailNow("Failed find user: ", err)
	}

	if res.ID.String() != user.ID.String() {
		s.FailNow("User id mismatch")
	}
	if req.Username != user.Username {
		s.FailNow("User username mismatch")
	}
}

func (s *AuthServiceTestSuite) SetupSuite() {
	s.PostgresContainerSuite.SetupSuite()
	logger := log.NewTestLogger(s.T())
	ctxManager := db.NewSqlContextManager(logger, s.DB)
	userRepository := repository.NewPgUserRepository(logger)
	validator := app.NewGoPlaygroundValidator(logger)
	passwordHasher := utils.NewBcryptPasswordHasher()
	authService := service.NewAuthServiceImpl(logger, ctxManager, userRepository, validator, passwordHasher)

	s.userRepository = userRepository
	s.authService = authService
	s.ctxManager = ctxManager
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}
