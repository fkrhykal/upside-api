package service

import (
	"context"

	"github.com/fkrhykal/upside-api/internal/account/dto"
	"github.com/fkrhykal/upside-api/internal/account/entity"
	"github.com/fkrhykal/upside-api/internal/account/repository"
	"github.com/fkrhykal/upside-api/internal/account/utils"
	"github.com/fkrhykal/upside-api/internal/shared/db"
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
) AuthService {
	return &AuthServiceImpl[T]{
		logger:         logger,
		userRepository: userRepository,
		ctxManager:     ctxManager,
		validator:      validator,
		passwordHasher: passwordHasher,
	}
}

type AuthServiceImpl[T any] struct {
	logger         log.Logger
	userRepository repository.UserRepository[T]
	ctxManager     db.CtxManager[T]
	validator      validation.Validator
	passwordHasher utils.PasswordHasher
}

func (s AuthServiceImpl[T]) SignUp(ctx context.Context, request *dto.SignUpRequest) (*dto.SignUpResponse, error) {
	s.logger.Infof("Attempting to register user with username: %s", request.Username)

	err := s.validator.Validate(request)

	validationError, ok := err.(*validation.ValidationError)
	if !ok {
		s.logger.Errorf("Failed to register user caused by: %+v", err)
		return nil, err
	}
	if validationError.Exist("username") {
		s.logger.Warnf("User registration failed due to validation error: %+v", validationError)
		return nil, err
	}

	dbCtx := s.ctxManager.NewDBContext(ctx)

	user, err := s.userRepository.FindByUsername(dbCtx, request.Username)
	if err != nil {
		s.logger.Errorf("Failed to register user caused by: %+v", err)
		return nil, err
	}
	if user != nil {
		validationError.Add("username", "username already used")
		s.logger.Warnf("Username already exists: %s, registration failed", request.Username)
		return nil, validationError
	}

	hashedPassword, err := s.passwordHasher.Hash(request.Password)
	if err != nil {
		s.logger.Errorf("Failed to register user due to password hashing failure: %+v", err)
		return nil, err
	}

	user = &entity.User{
		ID:       uuid.New(),
		Username: request.Username,
		Password: hashedPassword,
	}

	s.logger.Debugf("Preparing to save new user: %s", user.Username) // Detailed info for debugging
	if err := s.userRepository.Save(dbCtx, user); err != nil {
		s.logger.Errorf("Failed to save user %s: %v", user.Username, err)
		return nil, err
	}

	s.logger.Infof("Successfully registered user: %s", user.Username)

	return &dto.SignUpResponse{ID: user.ID}, nil
}
