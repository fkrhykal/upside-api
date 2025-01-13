package db

import (
	"context"
)

type DBContext[T any] interface {
	context.Context
	Executor() T
}

type TxContext[T any] interface {
	DBContext[T]
	Rollback() error
	Commit() error
}
