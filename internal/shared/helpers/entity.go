package helpers

import (
	"time"

	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

func SetupUser(ctx db.DBContext[db.SqlExecutor], s *suite.Suite) uuid.UUID {
	query := `INSERT INTO users(id, username, password) VALUES($1, $2, $3)`
	id := uuid.New()
	username := faker.Username()
	password := faker.Password()
	_, err := ctx.Executor().ExecContext(ctx, query, id, username, password)
	s.Require().NoError(err)
	return id
}

func SetupSide(ctx db.DBContext[db.SqlExecutor], s *suite.Suite) uuid.UUID {
	query := `INSERT INTO sides(id, nick, name, description, created_at) VALUES($1, $2, $3, $4, $5)`
	id := uuid.New()
	nick := faker.Username()
	name := faker.Name()
	description := faker.Sentence()
	createdAt := time.Now().UnixMilli()

	_, err := ctx.Executor().ExecContext(ctx, query, id, nick, name, description, createdAt)
	s.Require().NoError(err)

	return id
}

func SetupMembership(ctx db.DBContext[db.SqlExecutor], s *suite.Suite, userID uuid.UUID, sideID uuid.UUID) uuid.UUID {
	query := `INSERT INTO memberships(id, member_id, side_id, role) VALUES($1, $2, $3, $4)`
	id := uuid.New()
	_, err := ctx.Executor().ExecContext(ctx, query, id, userID, sideID, entity.MEMBER)
	s.Require().NoError(err)
	return id
}
