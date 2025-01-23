package service

import (
	"context"

	"github.com/fkrhykal/upside-api/internal/account/dto"
	"github.com/fkrhykal/upside-api/internal/account/repository"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/exception"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/google/uuid"
)

type UserServiceImpl[T any] struct {
	logger         log.Logger
	ctxManager     db.CtxManager[T]
	userRepository repository.UserRepository[T]
}

func NewUserServiceImpl[T any](logger log.Logger, ctxManger db.CtxManager[T], userRepository repository.UserRepository[T]) UserService {
	return &UserServiceImpl[T]{
		logger:         logger,
		ctxManager:     ctxManger,
		userRepository: userRepository,
	}
}

func (us *UserServiceImpl[T]) GetUserDetail(ctx context.Context, id uuid.UUID) (*dto.UserDetail, error) {
	dbCtx := us.ctxManager.NewDBContext(ctx)
	user, err := us.userRepository.FindById(dbCtx, id)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return &dto.UserDetail{ID: user.ID, Username: user.Username}, nil
	}
	return nil, exception.UserNotFound
}
