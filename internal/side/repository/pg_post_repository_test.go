package repository_test

import (
	"context"
	"testing"

	"github.com/fkrhykal/upside-api/internal/app"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/helpers"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/pagination"
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

func (p *PgPostRepositorySuite) TestFindManyWithEmptyULIDCursor() {
	ctx := context.Background()
	dbCtx := p.ctxManager.NewDBContext(ctx)

	sideID := helpers.SetupSide(dbCtx, &p.Suite)
	userID := helpers.SetupUser(dbCtx, &p.Suite)
	postsIDs := helpers.SetupPosts(dbCtx, &p.Suite, userID, sideID, 10)

	limitCursor := &pagination.NextULIDCursor{}
	posts, err := p.postRepository.FindManyWithULIDCursor(dbCtx, limitCursor)
	p.Require().NoError(err)
	p.Require().Len(posts, 10)

	p.EqualValues(postsIDs[9], posts[0].ID)
}

func (p *PgPostRepositorySuite) TestFindManyWithLimitULIDCursor() {
	ctx := context.Background()
	dbCtx := p.ctxManager.NewDBContext(ctx)

	sideID := helpers.SetupSide(dbCtx, &p.Suite)
	userID := helpers.SetupUser(dbCtx, &p.Suite)
	postsIDs := helpers.SetupPosts(dbCtx, &p.Suite, userID, sideID, 10)

	limitCursor := pagination.LimitNextULIDCursor(2)
	posts, err := p.postRepository.FindManyWithULIDCursor(dbCtx, limitCursor)
	p.Require().NoError(err)
	p.Require().Len(posts, limitCursor.Limit()+1)

	p.EqualValues(postsIDs.Last(), posts[0].ID)
	p.EqualValues(postsIDs[postsIDs.LastIndex()-limitCursor.Limit()], posts.Last().ID)
}

func (p *PgPostRepositorySuite) TestFindManyWithULIDCursor() {
	ctx := context.Background()
	dbCtx := p.ctxManager.NewDBContext(ctx)

	sideID := helpers.SetupSide(dbCtx, &p.Suite)
	userID := helpers.SetupUser(dbCtx, &p.Suite)
	postsIDs := helpers.SetupPosts(dbCtx, &p.Suite, userID, sideID, 20)

	limitCursor := pagination.LimitNextULIDCursor(2)
	posts, err := p.postRepository.FindManyWithULIDCursor(dbCtx, limitCursor)
	p.Require().NoError(err)
	p.Require().Len(posts, 3)

	for i, post := range posts {
		p.Require().EqualValues(postsIDs[19-i], post.ID)
	}

	nextCursor, err := pagination.NewNextULIDCursor(&posts.Second().ID, limitCursor.Limit())
	p.Require().NoError(err)

	nextPosts, err := p.postRepository.FindManyWithULIDCursor(dbCtx, nextCursor)
	p.Require().NoError(err)
	p.Require().Len(nextPosts, 3)

	for i, id := range postsIDs {
		if nextCursor.ID().String() == id.String() {
			p.Debugf("Equal %d", i)
		}

		latestPostID := nextPosts[0].ID.String()
		oldestPostID := nextPosts.Last().ID.String()
		currentPostID := postsIDs[i].String()

		if latestPostID == currentPostID {
			p.Debugf("%d: lastedPostID: %s", i, latestPostID)
		}

		if oldestPostID == currentPostID {
			p.Debugf("%d: oldestPostID: %s", i, oldestPostID)
		}
	}
}

func (p *PgPostRepositorySuite) TestFindManyWithLimit() {
	ctx := context.Background()
	dbCtx := p.ctxManager.NewDBContext(ctx)

	sideID := helpers.SetupSide(dbCtx, &p.Suite)
	userID := helpers.SetupUser(dbCtx, &p.Suite)
	helpers.SetupPosts(dbCtx, &p.Suite, userID, sideID, 10)

	posts, err := p.postRepository.FindManyWithLimit(dbCtx, 5)
	p.Require().NoError(err)
	p.Require().Len(posts, 5)
}

func (p *PgPostRepositorySuite) TestSavePost() {
	ctx := context.Background()
	dbCtx := p.ctxManager.NewDBContext(ctx)

	sideID := helpers.SetupSide(dbCtx, &p.Suite)
	userID := helpers.SetupUser(dbCtx, &p.Suite)

	post := entity.NewPost(faker.Sentence(), userID, sideID)

	err := p.postRepository.Save(dbCtx, post)
	p.Require().NoError(err)
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
