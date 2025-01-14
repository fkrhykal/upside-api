package db

import (
	"context"
	"database/sql"

	"github.com/fkrhykal/upside-api/internal/shared/log"
)

type SqlContextManager struct {
	logger log.Logger
	db     *sql.DB
}

func NewSqlContextManager(logger log.Logger, db *sql.DB) CtxManager[SqlExecutor] {
	return &SqlContextManager{
		logger: logger,
		db:     db,
	}
}

func (ctxManager *SqlContextManager) NewDBContext(ctx context.Context) DBContext[SqlExecutor] {
	return &SqlDBContext{Context: ctx, executor: ctxManager.db}
}

func (ctxManager *SqlContextManager) NewTxContext(ctx context.Context) (TxContext[SqlExecutor], error) {
	tx, err := ctxManager.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	return &SqlTxContext{Context: ctx, executor: tx}, nil
}
