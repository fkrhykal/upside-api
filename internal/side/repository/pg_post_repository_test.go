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
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/suite"
)

type PgPostRepositorySuite struct {
	app.PostgresContainerSuite
	postRepository repository.PostRepository[db.SqlExecutor]
	ctxManager     db.CtxManager[db.SqlExecutor]
}

func (p *PgPostRepositorySuite) TestSavePost() {
	ctx := context.Background()
	dbCtx := p.ctxManager.NewDBContext(ctx)

	sideID := p.setupSide(dbCtx)
	userID := p.setupUser(dbCtx)

	post := &entity.Post{
		ID:        ulid.Make(),
		Body:      faker.Sentence(),
		CreatedAt: time.Now().UnixMilli(),
		Author:    &entity.Author{ID: userID},
		Side:      &entity.Side{ID: sideID},
	}

	err := p.postRepository.Save(dbCtx, post)
	p.Require().NoError(err)
}

func (p *PgPostRepositorySuite) setupUser(ctx db.DBContext[db.SqlExecutor]) uuid.UUID {
	query := `INSERT INTO users(id, username, password) VALUES($1, $2, $3)`
	id := uuid.New()
	username := faker.Username()
	password := faker.Password()
	_, err := ctx.Executor().ExecContext(ctx, query, id, username, password)
	p.Require().NoError(err)
	return id
}

func (p *PgPostRepositorySuite) setupSide(ctx db.DBContext[db.SqlExecutor]) uuid.UUID {
	query := `INSERT INTO sides(id, nick, name, description, created_at) VALUES($1, $2, $3, $4, $5)`
	id := uuid.New()
	nick := faker.Username()
	name := faker.Name()
	description := faker.Sentence()
	createdAt := time.Now().UnixMilli()

	_, err := ctx.Executor().ExecContext(ctx, query, id, nick, name, description, createdAt)
	p.Require().NoError(err)

	return id
}

func (p *PgPostRepositorySuite) SetupSuite() {
	p.PostgresContainerSuite.SetupSuite()
	logger := log.NewTestLogger(p.T())
	p.ctxManager = db.NewSqlContextManager(logger, p.DB)
	p.postRepository = repository.NewPgPostRepository(logger)
}

func TestPgPostRepository(t *testing.T) {
	suite.Run(t, new(PgPostRepositorySuite))
}
