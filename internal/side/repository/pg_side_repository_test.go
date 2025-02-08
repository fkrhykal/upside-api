package repository_test

import (
	"context"
	"math/rand/v2"
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

func (mr *PgSideRepositorySuite) TestFindManyIn() {
	ids := make([]uuid.UUID, 100)
	ctx := context.Background()
	dbCtx := mr.ctxManager.NewDBContext(ctx)

	for i := range 100 {
		s := &entity.Side{
			ID:          uuid.New(),
			Nick:        faker.Username(),
			Name:        faker.Name(),
			Description: faker.Sentence(),
			CreatedAt:   uint64(time.Now().UnixMilli()),
		}

		err := mr.sideRepository.Save(dbCtx, s)
		mr.Nil(err)
		ids[i] = s.ID
	}

	sides, err := mr.sideRepository.FindManyIn(dbCtx, ids)
	mr.Nil(err)
	mr.Len(sides, 100)

	for i, side := range sides {
		mr.EqualValues(ids[i], side.ID)
	}
}

func (mr *PgSideRepositorySuite) TestFindLimitedWithLargestMemberships() {
	sideIDs := make([]uuid.UUID, 10)
	ctx := context.Background()
	dbCtx := mr.ctxManager.NewDBContext(ctx)

	for i := range 10 {
		s := &entity.Side{
			ID:          uuid.New(),
			Nick:        faker.Username(),
			Name:        faker.Name(),
			Description: faker.Sentence(),
			CreatedAt:   uint64(time.Now().UnixMilli()),
		}

		err := mr.sideRepository.Save(dbCtx, s)
		mr.Nil(err)
		sideIDs[i] = s.ID
	}

	userIDs := make([]uuid.UUID, 1000)

	insertUserQuery := `INSERT INTO users(id, username, password) VALUES($1, $2, $3)`
	insertUserStmt, err := dbCtx.Executor().PrepareContext(ctx, insertUserQuery)
	mr.Nil(err)

	for i := range 1000 {
		id := uuid.New()
		username := faker.Username()
		password := faker.Password()
		_, err := insertUserStmt.ExecContext(ctx, id, username, password)
		mr.Nil(err)
		userIDs[i] = id
	}

	insertMembershipQuery := `INSERT INTO memberships(id, member_id, side_id, role) VALUES($1, $2, $3, $4)`
	insertMembershipStmt, err := dbCtx.Executor().PrepareContext(ctx, insertMembershipQuery)
	mr.Nil(err)

	for _, userID := range userIDs {
		memberID := uuid.New()
		_, err = insertMembershipStmt.ExecContext(ctx, memberID, userID, sideIDs[rand.IntN(10)], entity.MEMBER)
		mr.Nil(err)
	}

	sides, err := mr.sideRepository.FindLimitedWithLargestMemberships(dbCtx, 5)
	mr.Nil(err)
	mr.Len(sides, 5)
}

func (mr *PgSideRepositorySuite) TestFindOffsetLimitedWithLargestMemberships() {
	sideIDs := make([]uuid.UUID, 10)
	ctx := context.Background()
	dbCtx := mr.ctxManager.NewDBContext(ctx)

	for i := range 10 {
		s := &entity.Side{
			ID:          uuid.New(),
			Nick:        faker.Username(),
			Name:        faker.Name(),
			Description: faker.Sentence(),
			CreatedAt:   uint64(time.Now().UnixMilli()),
		}

		err := mr.sideRepository.Save(dbCtx, s)
		mr.Nil(err)
		sideIDs[i] = s.ID
	}

	userIDs := make([]uuid.UUID, 1000)

	insertUserQuery := `INSERT INTO users(id, username, password) VALUES($1, $2, $3)`
	insertUserStmt, err := dbCtx.Executor().PrepareContext(ctx, insertUserQuery)
	mr.Nil(err)

	for i := range 1000 {
		id := uuid.New()
		username := faker.Username()
		password := faker.Password()
		_, err := insertUserStmt.ExecContext(ctx, id, username, password)
		mr.Nil(err)
		userIDs[i] = id
	}

	insertMembershipQuery := `INSERT INTO memberships(id, member_id, side_id, role) VALUES($1, $2, $3, $4)`
	insertMembershipStmt, err := dbCtx.Executor().PrepareContext(ctx, insertMembershipQuery)
	mr.Nil(err)

	for _, userID := range userIDs {
		memberID := uuid.New()
		_, err = insertMembershipStmt.ExecContext(ctx, memberID, userID, sideIDs[rand.IntN(10)], entity.MEMBER)
		mr.Nil(err)
	}

	sides, err := mr.sideRepository.FindOffsetLimitedWithLargestMemberships(dbCtx, 1, 5)
	mr.Nil(err)
	mr.Len(sides, 5)
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
