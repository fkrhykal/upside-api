package repository

import (
	"github.com/fkrhykal/upside-api/internal/account/entity"
	"github.com/fkrhykal/upside-api/internal/shared/db"
)

type PgUserRepository[T db.SqlExecutor] struct {
}

func NewPgUserRepository() UserRepository[db.SqlExecutor] {
	return &PgUserRepository[db.SqlExecutor]{}
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

	if err != nil {
		return nil, err
	}
	return user, nil
}
