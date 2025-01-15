package repository

import (
	"database/sql"
	"errors"

	"github.com/fkrhykal/upside-api/internal/account/entity"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/log"
)

type PgUserRepository[T db.SqlExecutor] struct {
	logger log.Logger
}

func NewPgUserRepository(logger log.Logger) UserRepository[db.SqlExecutor] {
	return &PgUserRepository[db.SqlExecutor]{
		logger,
	}
}

func (r *PgUserRepository[T]) Save(ctx db.DBContext[T], user *entity.User) error {
	query := "INSERT INTO users(id, username, password) VALUES($1,$2,$3)"
	_, err := ctx.Executor().
		ExecContext(ctx, query, user.ID, user.Username, user.Password)
	return err
}

func (r *PgUserRepository[T]) FindByUsername(ctx db.DBContext[T], username string) (*entity.User, error) {
	user := new(entity.User)
	query := "SELECT id, username, password FROM users WHERE username = $1"
	err := ctx.Executor().
		QueryRowContext(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
