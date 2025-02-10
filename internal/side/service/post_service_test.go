package service_test

import (
	"context"
	"testing"

	"github.com/fkrhykal/upside-api/internal/app"
	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/helpers"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/side/dto"
	"github.com/fkrhykal/upside-api/internal/side/repository"
	"github.com/fkrhykal/upside-api/internal/side/service"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/suite"
)

type PostServiceSuite struct {
	app.PostgresContainerSuite
	postService          service.PostService
	ctxManager           db.CtxManager[db.SqlExecutor]
	postRepository       repository.PostRepository[db.SqlExecutor]
	membershipRepository repository.MembershipRepository[db.SqlExecutor]
	sideRepository       repository.SideRepository[db.SqlExecutor]
}

func (ps *PostServiceSuite) TestCreateSide() {
	ctx := context.Background()
	dbCtx := ps.ctxManager.NewDBContext(ctx)

	sideID := helpers.SetupSide(dbCtx, &ps.Suite)
	userID := helpers.SetupUser(dbCtx, &ps.Suite)
	helpers.SetupMembership(dbCtx, &ps.Suite, userID, sideID)

	authCtx := auth.NewAuthContext(ctx, &auth.UserCredential{ID: userID})

	res, err := ps.postService.CreatePost(authCtx, &dto.CreatePostRequest{
		SideID: sideID,
		Body:   faker.Sentence(),
	})
	ps.Require().NoError(err)
	ps.Require().NotNil(res)
}

func (ps *PostServiceSuite) SetupSuite() {
	ps.PostgresContainerSuite.SetupSuite()
	logger := log.NewTestLogger(ps.T())

	ps.membershipRepository = repository.NewPgMembershipRepository(logger)
	ps.postRepository = repository.NewPgPostRepository(logger)
	ps.sideRepository = repository.NewPgSideRepository(logger)
	ps.ctxManager = db.NewSqlContextManager(logger, ps.DB)

	validator := app.NewGoPlaygroundValidator(logger)

	ps.postService = service.NewPostServiceImpl(
		validator,
		ps.ctxManager,
		ps.sideRepository,
		ps.membershipRepository,
		ps.postRepository,
	)
}

func TestPostService(t *testing.T) {
	suite.Run(t, new(PostServiceSuite))
}
