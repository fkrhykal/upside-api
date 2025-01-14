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
	s.logger.Infof("Received sign-up request for username: %s", request.Username)

	if err := s.validator.Validate(request); err != nil {
		s.logger.Warnf("Validation failed: %v", err)
		return nil, err
	}

	dbCtx := s.ctxManager.NewDBContext(ctx)

	hashedPassword, err := s.passwordHasher.Hash(request.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:       uuid.New(),
		Username: request.Username,
		Password: hashedPassword,
	}

	s.logger.Debugf("Attempting to save user: %s", user.Username)
	if err := s.userRepository.Save(dbCtx, user); err != nil {
		s.logger.Errorf("Failed to save user %s: %v", user.Username, err)
		return nil, err
	}

	s.logger.Infof("Successfully registered user: %s", user.Username)
	return &dto.SignUpResponse{ID: user.ID}, nil
}
