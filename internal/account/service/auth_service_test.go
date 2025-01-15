package service_test

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/suite"

	"github.com/fkrhykal/upside-api/internal/account/dto"
	"github.com/fkrhykal/upside-api/internal/account/repository"
	"github.com/fkrhykal/upside-api/internal/account/service"
	"github.com/fkrhykal/upside-api/internal/account/utils"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	s "github.com/fkrhykal/upside-api/internal/shared/suite"
	"github.com/fkrhykal/upside-api/internal/shared/validation"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type AuthServiceTestSuite struct {
	s.PostgresContainerSuite
	authService    service.AuthService
	userRepository repository.UserRepository[db.SqlExecutor]
	ctxManager     db.CtxManager[db.SqlExecutor]
}

func (s *AuthServiceTestSuite) TestSignUp() {
	t := s.T()

	req := &dto.SignUpRequest{
		Username: faker.Username(),
		Password: faker.Password(),
	}

	res, err := s.authService.SignUp(s.Ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("sign-up response: %+v", res)

	dbCtx := s.ctxManager.NewDBContext(s.Ctx)
	user, err := s.userRepository.FindByUsername(dbCtx, req.Username)
	if err != nil {
		t.Fatal(err)
	}

	if res.ID.String() != user.ID.String() {
		t.Fatal("user id mismatch")
	}
	if req.Username != user.Username {
		t.Fatal("user username mismatch")
	}
	t.Logf("user: %+v \n", user)
}

func (s *AuthServiceTestSuite) SetupSuite() {
	s.PostgresContainerSuite.SetupSuite()

	t := s.T()
	logger := log.NewTestLogger(t)
	ctxManager := db.NewSqlContextManager(logger, s.DB)
	userRepository := repository.NewPgUserRepository(logger)
	validator := validation.NewGoPlaygroundValidator(logger)
	passwordHasher := utils.NewBcryptPasswordHasher()
	authService := service.NewAuthServiceImpl(logger, ctxManager, userRepository, validator, passwordHasher)

	s.userRepository = userRepository
	s.authService = authService
	s.ctxManager = ctxManager
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}
