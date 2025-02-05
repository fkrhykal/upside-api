package db

type TxContext[T any] interface {
	DBContext[T]
	Rollback() error
	Commit() error
}
