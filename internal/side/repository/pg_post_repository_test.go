package repository_test

import (
	"context"
	"testing"

	"github.com/fkrhykal/upside-api/internal/app"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/helpers"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/side/entity"
	"github.com/fkrhykal/upside-api/internal/side/repository"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/suite"
)

type PgPostRepositorySuite struct {
	app.PostgresContainerSuite
	log.Logger
	postRepository repository.PostRepository[db.SqlExecutor]
	ctxManager     db.CtxManager[db.SqlExecutor]
}

func (p *PgPostRepositorySuite) TestSavePost() {
	ctx := context.Background()
	dbCtx := p.ctxManager.NewDBContext(ctx)

	sideID := helpers.SetupSide(dbCtx, &p.Suite)
	userID := helpers.SetupUser(dbCtx, &p.Suite)

	post := entity.CreatePost(faker.Sentence(), userID, sideID)

	err := p.postRepository.Save(dbCtx, post)
	p.Require().NoError(err)

}

func (p *PgPostRepositorySuite) TestFindByID() {
	ctx := context.Background()
	dbCtx := p.ctxManager.NewDBContext(ctx)

	sideID := helpers.SetupSide(dbCtx, &p.Suite)
	userID := helpers.SetupUser(dbCtx, &p.Suite)
	postIDs := helpers.SetupPosts(dbCtx, &p.Suite, userID, sideID, 1)

	post, err := p.postRepository.FindByID(dbCtx, postIDs[0])
	p.Require().NoError(err)
	p.Require().EqualValues(postIDs[0], post.ID)
	p.Require().EqualValues(sideID, post.Side.ID)
	p.Require().EqualValues(userID, post.Author.ID)
}

func (p *PgPostRepositorySuite) SetupSuite() {
	p.PostgresContainerSuite.SetupSuite()
	p.Logger = log.NewTestLogger(p.T())
	p.ctxManager = db.NewSqlContextManager(p.Logger, p.DB)
	p.postRepository = repository.NewPgPostRepository(p.Logger)
}

func TestPgPostRepository(t *testing.T) {
	suite.Run(t, new(PgPostRepositorySuite))
}
