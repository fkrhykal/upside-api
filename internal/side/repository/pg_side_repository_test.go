package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/fkrhykal/upside-api/internal/app"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/fkrhykal/upside-api/internal/side/repository"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type PgSideRepositorySuite struct {
	app.PostgresContainerSuite
	sideRepository repository.SideRepository[db.SqlExecutor]
	ctxManager     db.CtxManager[db.SqlExecutor]
}

func (mr *PgSideRepositorySuite) TestSaveSide() {
	ctx := context.Background()
	dbCtx := mr.ctxManager.NewDBContext(ctx)

	side := &entity.Side{
		ID:          uuid.New(),
		Nick:        faker.Username(),
		Name:        faker.Name(),
		Description: faker.Sentence(),
		CreatedAt:   uint64(time.Now().UnixMilli()),
	}

	err := mr.sideRepository.Save(dbCtx, side)
	mr.Nil(err)

	var id uuid.UUID
	err = mr.DB.QueryRowContext(ctx, "SELECT id FROM sides WHERE id=$1", side.ID).Scan(&id)
	mr.Nil(err)
	mr.EqualValues(side.ID, id)
}

func (mr *PgSideRepositorySuite) TestFindByID() {
	ctx := context.Background()
	dbCtx := mr.ctxManager.NewDBContext(ctx)

	s := &entity.Side{
		ID:          uuid.New(),
		Nick:        faker.Username(),
		Name:        faker.Name(),
		Description: faker.Sentence(),
		CreatedAt:   uint64(time.Now().UnixMilli()),
	}

	err := mr.sideRepository.Save(dbCtx, s)
	mr.Nil(err)

	side, err := mr.sideRepository.FindById(dbCtx, s.ID)
	mr.Nil(err)
	mr.T().Log(side.CreatedAt)
	mr.EqualValues(s.CreatedAt, side.CreatedAt)
}

func (mr *PgSideRepositorySuite) SetupSuite() {
	mr.PostgresContainerSuite.SetupSuite()
	logger := log.NewTestLogger(mr.T())
	mr.ctxManager = db.NewSqlContextManager(logger, mr.DB)
	mr.sideRepository = repository.NewPgSideRepository(logger)
}

func TestPgSideRepository(t *testing.T) {
	suite.Run(t, new(PgSideRepositorySuite))
}
