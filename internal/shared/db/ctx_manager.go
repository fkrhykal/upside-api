package db

import "context"

type CtxManager[T any] interface {
	NewTxContext(context.Context) (TxContext[T], error)
	NewDBContext(context.Context) DBContext[T]
}
			