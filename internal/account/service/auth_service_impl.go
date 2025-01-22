package service

import (
	"context"
	"time"

	"github.com/fkrhykal/upside-api/internal/account/dto"
	"github.com/fkrhykal/upside-api/internal/account/entity"
	"github.com/fkrhykal/upside-api/internal/account/repository"
	"github.com/fkrhykal/upside-api/internal/account/utils"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/exception"
	"github.com/fkrhykal/upside-api/internal/shared/helpers"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/validation"

	"github.com/google/uuid"
)

func NewAuthServiceImpl[T any](
	logger log.Logger,
	ctxManager db.CtxManager[T],
	userRepository repository.UserRepository[T],
	validator validation.Validator,
	passwordHasher utils.PasswordHasher,
	credentialService CredentialService,
) AuthService {
	return &AuthServiceImpl[T]{
		logger:            logger,
		userRepository:    userRepository,
		ctxManager:        ctxManager,
		validator:         validator,
		passwordHasher:    passwordHasher,
		credentialService: credentialService,
	}
}

type AuthServiceImpl[T any] struct {
	logger            log.Logger
	userRepository    repository.UserRepository[T]
	ctxManager        db.CtxManager[T]
	validator         validation.Validator
	passwordHasher    utils.PasswordHasher
	credentialService CredentialService
}

func (s *AuthServiceImpl[T]) SignUp(ctx context.Context, req *dto.SignUpRequest) (*dto.SignUpResponse, error) {
	s.logger.Infof("Attempting to register user with username: %s", req.Username)

	err := s.validator.Validate(req)
	if err != nil {
		return nil, err
	}

	dbCtx := s.ctxManager.NewDBContext(ctx)

	user, err := s.userRepository.FindByUsername(dbCtx, req.Username)
	if err != nil {
		s.logger.Errorf("Failed to register user caused by: %+v", err)
		return nil, err
	}
	if user != nil {
		s.logger.Warnf("Username already exists: %s, registration failed", req.Username)
		return nil, &validation.ValidationError{
			Detail: validation.ErrorDetail{
				"username": "username already used",
			},
		}
	}

	hashedPassword, err := s.passwordHasher.Hash(req.Password)
	if err != nil {
		s.logger.Errorf("Failed to register user due to password hashing failure: %+v", err)
		return nil, err
	}

	user = &entity.User{
		ID:       uuid.New(),
		Username: req.Username,
		Password: hashedPassword,
	}

	s.logger.Debugf("Preparing to save new user: %s", user.Username)
	if err := s.userRepository.Save(dbCtx, user); err != nil {
		s.logger.Errorf("Failed to save user %s: %v", user.Username, err)
		return nil, err
	}

	s.logger.Infof("Successfully registered user: %s", user.Username)

	return &dto.SignUpResponse{ID: user.ID}, nil
}

func (s *AuthServiceImpl[T]) SignIn(ctx context.Context, req *dto.SignInRequest) (*dto.SignInResponse, error) {
	s.logger.Infof("SignIn attempt started for username: %s", req.Username)

	err := s.validator.Validate(req)
	if err != nil {
		s.logger.Errorf("Validation failed for SignIn request: %+v", err)
		return nil, exception.ErrAuthentication
	}

	dbCtx := s.ctxManager.NewDBContext(ctx)
	s.logger.Debugf("Attempting to retrieve user with username: %s", req.Username)

	user, err := s.userRepository.FindByUsername(dbCtx, req.Username)
	if err != nil {
		s.logger.Errorf("Error while retrieving user by username %s: %v", req.Username, err)
		return nil, err
	}

	if user == nil {
		s.logger.Warnf("Authentication failed: no user found with username %s", req.Username)
		return nil, exception.ErrAuthentication
	}

	s.logger.Debug("Matching provided password with stored hash")
	if matches := s.passwordHasher.Match(req.Password, user.Password); !matches {
		s.logger.Warnf("Authentication failed: password mismatch for username %s", req.Username)
		return nil, exception.ErrAuthentication
	}

	s.logger.Debugf("Generating token for user ID: %s", user.ID)
	token, err := s.credentialService.GenerateToken(ctx, &dto.UserCredential{ID: user.ID}, time.Now().Add(helpers.WEEK))
	if err != nil {
		s.logger.Errorf("Error generating token for user ID %s: %v", user.ID, err)
		return nil, err
	}

	s.logger.Infof("User signed in successfully: %s", req.Username)
	return &dto.SignInResponse{Token: token}, nil
}
