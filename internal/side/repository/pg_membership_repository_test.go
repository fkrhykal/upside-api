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

type PgMembershipRepositorySuite struct {
	app.PostgresContainerSuite
	membershipRepository repository.MembershipRepository[db.SqlExecutor]
	ctxManager           db.CtxManager[db.SqlExecutor]
}

func (mr *PgMembershipRepositorySuite) TestSaveMembership() {
	ctx := context.Background()
	dbCtx := mr.ctxManager.NewDBContext(ctx)

	memberId := mr.saveUser(dbCtx)
	sideId := mr.saveSide(dbCtx)

	membership := &entity.Membership{
		ID:     uuid.New(),
		Member: memberId,
		Side:   sideId,
		Role:   entity.FOUNDER,
	}

	err := mr.membershipRepository.Save(dbCtx, membership)
	mr.Nil(err)

	var id uuid.UUID
	err = mr.DB.QueryRowContext(ctx, "SELECT id FROM memberships WHERE id=$1", membership.ID).Scan(&id)
	mr.Nil(err)
	mr.EqualValues(membership.ID, id)
}

func (mr *PgMembershipRepositorySuite) TestGetJoinedSides() {
	ctx := context.Background()
	dbCtx := mr.ctxManager.NewDBContext(ctx)

	memberId := mr.saveUser(dbCtx)
	sideId := mr.saveMultipleSide(dbCtx, 2)

	membershipRegistry := make([]*entity.Membership, 100)

	for i := range 50 {
		membership := &entity.Membership{
			ID:     uuid.New(),
			Member: memberId,
			Side:   sideId[0],
			Role:   entity.FOUNDER,
		}

		err := mr.membershipRepository.Save(dbCtx, membership)
		mr.Nil(err)

		membershipRegistry[i] = membership
	}

	for i := range 50 {
		membership := &entity.Membership{
			ID:     uuid.New(),
			Member: memberId,
			Side:   sideId[1],
			Role:   entity.FOUNDER,
		}

		err := mr.membershipRepository.Save(dbCtx, membership)
		mr.Nil(err)

		membershipRegistry[i+50] = membership
	}

	memberships, err := mr.membershipRepository.FindManyByMemberID(dbCtx, memberId)
	mr.Nil(err)
	for i, membership := range memberships {
		mr.EqualValues(membershipRegistry[i], membership)
	}
}

func (mr *PgMembershipRepositorySuite) TestFindManyBySideIDsAndMemberID() {
	ctx := context.Background()
	dbCtx := mr.ctxManager.NewDBContext(ctx)

	sideIDs := mr.saveMultipleSide(dbCtx, 10)
	userID := mr.saveUser(dbCtx)

	for _, sideID := range sideIDs {
		membership := &entity.Membership{
			ID:     uuid.New(),
			Member: userID,
			Side:   sideID,
			Role:   entity.MEMBER,
		}
		err := mr.membershipRepository.Save(dbCtx, membership)
		mr.Nil(err)
	}

	sides, err := mr.membershipRepository.FindManyBySideIDsAndMemberID(dbCtx, sideIDs, userID)
	mr.Nil(err)
	mr.Len(sides, 10)
}

func (mr *PgMembershipRepositorySuite) saveMultipleSide(dbCtx db.DBContext[db.SqlExecutor], amount int) uuid.UUIDs {
	IDs := make(uuid.UUIDs, amount)
	query := `INSERT INTO sides(id, nick, name, description, created_at) VALUES($1, $2, $3, $4, $5)`
	stmt, err := dbCtx.Executor().PrepareContext(dbCtx, query)
	mr.Nil(err)

	for i := range amount {
		id := uuid.New()
		nick := faker.Username()
		name := faker.Name()
		description := faker.Sentence()
		createdAt := time.Now().UnixMilli()
		_, err := stmt.ExecContext(dbCtx, id, nick, name, description, createdAt)
		mr.Nil(err)
		IDs[i] = id
	}

	return IDs
}

func (mr *PgMembershipRepositorySuite) saveSide(dbCtx db.DBContext[db.SqlExecutor]) uuid.UUID {
	id := uuid.New()
	nick := faker.Username()
	name := faker.Name()
	description := faker.Sentence()
	createdAt := time.Now().UnixMilli()
	query := "INSERT INTO sides(id, nick, name, description, created_at) VALUES($1,$2,$3,$4,$5)"
	_, err := dbCtx.Executor().
		ExecContext(dbCtx, query, id, nick, name, description, createdAt)
	mr.Nil(err)
	return id
}

func (mr *PgMembershipRepositorySuite) saveUser(dbCtx db.DBContext[db.SqlExecutor]) uuid.UUID {
	id := uuid.New()
	username := faker.Username()
	password := faker.Password()
	query := "INSERT INTO users(id, username, password) VALUES($1,$2,$3)"
	_, err := dbCtx.Executor().
		ExecContext(dbCtx, query, id, username, password)
	mr.Nil(err)
	return id
}

func (mr *PgMembershipRepositorySuite) SetupSuite() {
	mr.PostgresContainerSuite.SetupSuite()
	logger := log.NewTestLogger(mr.T())
	mr.ctxManager = db.NewSqlContextManager(logger, mr.DB)
	mr.membershipRepository = repository.NewPgMembershipRepository(logger)
}

func TestPgMembershipRepository(t *testing.T) {
	suite.Run(t, new(PgMembershipRepositorySuite))
}
