package service_test

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/fkrhykal/upside-api/internal/account/dto"
	"github.com/fkrhykal/upside-api/internal/account/entity"
	"github.com/fkrhykal/upside-api/internal/account/repository"
	"github.com/fkrhykal/upside-api/internal/account/service"
	"github.com/fkrhykal/upside-api/internal/account/utils"
	"github.com/fkrhykal/upside-api/internal/app"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/exception"
	"github.com/fkrhykal/upside-api/internal/shared/log"
)

type AuthServiceTestSuite struct {
	app.PostgresContainerSuite
	authService       service.AuthService
	userRepository    repository.UserRepository[db.SqlExecutor]
	ctxManager        db.CtxManager[db.SqlExecutor]
	passwordHasher    utils.PasswordHasher
	credentialService service.CredentialService
}

func (s *AuthServiceTestSuite) TestSignUp() {

	req := &dto.SignUpRequest{
		Username: faker.Username(),
		Password: "TestPassword123^$&",
	}

	ctx := context.Background()

	res, err := s.authService.SignUp(ctx, req)
	s.Nil(err, "Failed to sign up: %+v", err)

	dbCtx := s.ctxManager.NewDBContext(ctx)
	user, err := s.userRepository.FindByUsername(dbCtx, req.Username)
	s.Nil(err, "Failed to find user: %+v", err)

	s.EqualValues(res.ID, user.ID, "User id mismatch")
	s.Equal(req.Username, user.Username, "User username mismatch")
	if req.Username != user.Username {
		s.FailNow("User username mismatch")
	}
}

func (s *AuthServiceTestSuite) TestSignIn() {
	ctx := context.Background()

	password := "_%TestPassword123_"
	hashedPassword, err := s.passwordHasher.Hash(password)
	s.Nil(err, "Failed to hash password: %+v", err)

	user := &entity.User{
		ID:       uuid.New(),
		Username: faker.Username(),
		Password: hashedPassword,
	}

	dbCtx := s.ctxManager.NewDBContext(ctx)
	err = s.userRepository.Save(dbCtx, user)
	s.Nil(err, "Failed to save user: %+v", err)

	req := &dto.SignInRequest{
		Username: user.Username,
		Password: password,
	}

	res, err := s.authService.SignIn(ctx, req)
	s.Nil(err, "Failed to sign-in: %+v", err)

	credential, err := s.credentialService.RetrieveUserCredential(ctx, res.Token)
	s.Nil(err, "Failed to retrieve credential: %+v", err)

	s.EqualValues(user.ID, credential.ID, "User id and credential id mismatch: %+v", err)
}

func (s *AuthServiceTestSuite) TestSignInWrongPassword() {
	ctx := context.Background()

	password := "_%TestPassword123_"
	hashedPassword, err := s.passwordHasher.Hash(password)
	s.Nil(err, "Failed to hash password: %+v", err)

	user := &entity.User{
		ID:       uuid.New(),
		Username: faker.Username(),
		Password: hashedPassword,
	}

	dbCtx := s.ctxManager.NewDBContext(ctx)
	err = s.userRepository.Save(dbCtx, user)
	s.Nil(err, "Failed to save user: %+v", err)

	req := &dto.SignInRequest{
		Username: user.Username,
		Password: password + "&",
	}

	_, err = s.authService.SignIn(ctx, req)
	s.ErrorIs(err, exception.ErrAuthentication)
}

func (s *AuthServiceTestSuite) TestSignInUserNotFound() {
	ctx := context.Background()

	password := "_%TestPassword123_"
	req := &dto.SignInRequest{
		Username: faker.Username(),
		Password: password,
	}

	_, err := s.authService.SignIn(ctx, req)
	s.ErrorIs(err, exception.ErrAuthentication)
}

func (s *AuthServiceTestSuite) SetupSuite() {
	s.PostgresContainerSuite.SetupSuite()
	logger := log.NewTestLogger(s.T())
	ctxManager := db.NewSqlContextManager(logger, s.DB)
	userRepository := repository.NewPgUserRepository(logger)
	validator := app.NewGoPlaygroundValidator(logger)
	passwordHasher := utils.NewBcryptPasswordHasher()
	credentialService := service.NewJwtCredentialService(logger, &service.JwtCredentialConfig{
		SignedKey: []byte("secret"),
	})
	authService := service.NewAuthServiceImpl(logger, ctxManager, userRepository, validator, passwordHasher, credentialService)

	s.passwordHasher = passwordHasher
	s.userRepository = userRepository
	s.authService = authService
	s.ctxManager = ctxManager
	s.credentialService = credentialService
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}
