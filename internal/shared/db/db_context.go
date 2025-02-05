package db

import "context"

type DBContext[T any] interface {
	context.Context
	Executor() T
}
