package db_test

import (
	"context"
	"testing"

	"github.com/fkrhykal/upside-api/internal/app"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type SqlContextManagerSuite struct {
	app.PostgresContainerSuite
	ctxManager db.CtxManager[db.SqlExecutor]
}

func (cm *SqlContextManagerSuite) TestCommitTransaction() {
	ctx := context.Background()

	createTableQuery := "CREATE TABLE IF NOT EXISTS tests (id UUID PRIMARY KEY, message TEXT NOT NULL)"
	_, err := cm.DB.ExecContext(ctx, createTableQuery)
	cm.Nil(err)

	id := uuid.New()
	message := faker.Sentence()
	saveQuery := "INSERT INTO tests(id, message) VALUES($1, $2)"
	_, err = cm.DB.ExecContext(ctx, saveQuery, id, message)
	cm.Nil(err)

	tx, err := cm.ctxManager.NewTxContext(ctx)
	cm.Nil(err)
	cm.IsType(&db.SqlTxContext{}, tx)

	updateQuery := "UPDATE tests SET message = $1 WHERE id = $2"
	result, err := tx.Executor().ExecContext(tx, updateQuery, "Updated", id)
	cm.Nil(err)
	affected, err := result.RowsAffected()
	cm.Nil(err)
	cm.EqualValues(1, affected)

	err = tx.Commit()
	cm.Nil(err)

	var latestMessage string
	selectQuery := "SELECT message FROM tests WHERE id = $1"
	err = cm.DB.QueryRowContext(ctx, selectQuery, id.String()).Scan(&latestMessage)
	cm.Nil(err)
	cm.EqualValues("Updated", latestMessage)
}

func (cm *SqlContextManagerSuite) TestRollbackTransaction() {
	ctx := context.Background()

	createTableQuery := "CREATE TABLE IF NOT EXISTS tests (id UUID PRIMARY KEY, message TEXT NOT NULL)"
	_, err := cm.DB.ExecContext(ctx, createTableQuery)
	cm.Nil(err)

	id := uuid.New()
	message := faker.Sentence()
	saveQuery := "INSERT INTO tests(id, message) VALUES($1, $2)"
	_, err = cm.DB.ExecContext(ctx, saveQuery, id, message)
	cm.Nil(err)

	tx, err := cm.ctxManager.NewTxContext(ctx)
	cm.Nil(err)
	cm.IsType(&db.SqlTxContext{}, tx)

	updateQuery := "UPDATE tests SET message = $1 WHERE id = $2"
	result, err := tx.Executor().ExecContext(tx, updateQuery, "Updated", id)
	cm.Nil(err)
	affected, err := result.RowsAffected()
	cm.Nil(err)
	cm.EqualValues(1, affected)

	err = tx.Rollback()
	cm.Nil(err)

	var latestMessage string
	selectQuery := "SELECT message FROM tests WHERE id = $1"
	err = cm.DB.QueryRowContext(ctx, selectQuery, id.String()).Scan(&latestMessage)
	cm.Nil(err)
	cm.EqualValues(message, latestMessage)
}

func (cm *SqlContextManagerSuite) SetupSuite() {
	cm.PostgresContainerSuite.SetupSuite()
	logger := log.NewTestLogger(cm.T())
	cm.ctxManager = db.NewSqlContextManager(logger, cm.DB)
}

func TestSqlCtxManager(t *testing.T) {
	suite.Run(t, new(SqlContextManagerSuite))
}
