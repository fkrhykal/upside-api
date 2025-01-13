package service

import (
	"context"

	"github.com/fkrhykal/upside-api/internal/account/dto"
	"github.com/fkrhykal/upside-api/internal/account/entity"
	"github.com/fkrhykal/upside-api/internal/account/repository"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/validation"

	"github.com/google/uuid"
)

func NewAuthServiceImpl[T any](
	ctxManager db.CtxManager[T],
	userRepository repository.UserRepository[T],
	validator validation.Validator,
) AuthService {
	return &AuthServiceImpl[T]{
		userRepository: userRepository,
		ctxManager:     ctxManager,
		validator:      validator,
	}
}

type AuthServiceImpl[T any] struct {
	userRepository repository.UserRepository[T]
	ctxManager     db.CtxManager[T]
	validator      validation.Validator
}

func (s AuthServiceImpl[T]) SignUp(ctx context.Context, request *dto.SignUpRequest) (*dto.SignUpResponse, error) {

	if err := s.validator.Validate(request); err != nil {
		return nil, err
	}

	dbCtx := s.ctxManager.NewDBContext(ctx)

	user := &entity.User{
		ID:       uuid.New(),
		Username: request.Username,
		Password: request.Password,
	}

	if err := s.userRepository.Save(dbCtx, user); err != nil {
		return nil, err
	}

	return &dto.SignUpResponse{ID: user.ID}, nil
}
