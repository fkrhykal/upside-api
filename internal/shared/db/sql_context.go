package db

import (
	"context"
	"database/sql"
)

type SqlExecutor interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type SqlTxExecutor interface {
	SqlExecutor
	Rollback() error
	Commit() error
}

type SqlDBContext struct {
	context.Context
	executor SqlExecutor
}

func (dbCtx *SqlDBContext) Executor() SqlExecutor {
	return dbCtx.executor
}

type SqlTxContext struct {
	context.Context
	executor SqlTxExecutor
}

func (txCtx *SqlTxContext) Executor() SqlExecutor {
	return txCtx.executor
}

func (txCtx *SqlTxContext) Rollback() error {
	return txCtx.executor.Rollback()
}

func (txCtx *SqlTxContext) Commit() error {
	return txCtx.executor.Commit()
}
