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
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type PgVoteRepositorySuite struct {
	app.PostgresContainerSuite
	logger         log.Logger
	ctxManager     db.CtxManager[db.SqlExecutor]
	voteRepository repository.VoteRepository[db.SqlExecutor]
}

func (vs *PgVoteRepositorySuite) TestSaveVote() {
	ctx := context.Background()
	dbCtx := vs.ctxManager.NewDBContext(ctx)

	userID := helpers.SetupUser(dbCtx, &vs.Suite)
	sideID := helpers.SetupSide(dbCtx, &vs.Suite)
	postIDs := helpers.SetupPosts(dbCtx, &vs.Suite, userID, sideID, 1)

	vote := &entity.Vote{
		ID: uuid.New(),
		Voter: &entity.Voter{
			ID: userID,
		},
		Post: &entity.Post{
			ID: postIDs[0],
		},
		Kind: entity.UpVote,
	}

	err := vs.voteRepository.Save(dbCtx, vote)
	vs.Require().NoError(err)
}

func (vs *PgVoteRepositorySuite) TestFindByVoterIDAndPostID() {
	ctx := context.Background()
	dbCtx := vs.ctxManager.NewDBContext(ctx)

	userID := helpers.SetupUser(dbCtx, &vs.Suite)
	sideID := helpers.SetupSide(dbCtx, &vs.Suite)
	postIDs := helpers.SetupPosts(dbCtx, &vs.Suite, userID, sideID, 1)

	v := &entity.Vote{
		ID: uuid.New(),
		Voter: &entity.Voter{
			ID: userID,
		},
		Post: &entity.Post{
			ID: postIDs[0],
		},
		Kind: entity.UpVote,
	}

	err := vs.voteRepository.Save(dbCtx, v)
	vs.Require().NoError(err)

	vote, err := vs.voteRepository.FindByVoterIDAndPostID(dbCtx, v.Voter.ID, v.Post.ID)
	vs.Require().NoError(err)

	vs.Require().EqualValues(v.ID, vote.ID)
}

func (vs *PgVoteRepositorySuite) TestSumVoteKindGroupByPostIDs() {
	ctx := context.Background()
	dbCtx := vs.ctxManager.NewDBContext(ctx)

	authorID := helpers.SetupUser(dbCtx, &vs.Suite)
	voterID := helpers.SetupUser(dbCtx, &vs.Suite)
	sideID := helpers.SetupSide(dbCtx, &vs.Suite)
	postIDs := helpers.SetupPosts(dbCtx, &vs.Suite, authorID, sideID, 2)
	helpers.SetupVote(dbCtx, &vs.Suite, voterID, postIDs[0], entity.UpVote)

	registry, err := vs.voteRepository.SumVoteKindGroupByPostIDs(dbCtx, postIDs)
	vs.Require().NoError(err)
	vs.Require().EqualValues(1, registry[postIDs[0]])
	vs.Require().EqualValues(0, registry[postIDs[1]])
}

func (vs *PgVoteRepositorySuite) TestUpdate() {
	ctx := context.Background()
	dbCtx := vs.ctxManager.NewDBContext(ctx)

	authorID := helpers.SetupUser(dbCtx, &vs.Suite)
	voterID := helpers.SetupUser(dbCtx, &vs.Suite)
	sideID := helpers.SetupSide(dbCtx, &vs.Suite)
	postIDs := helpers.SetupPosts(dbCtx, &vs.Suite, authorID, sideID, 1)
	voteID := helpers.SetupVote(dbCtx, &vs.Suite, voterID, postIDs[0], entity.UpVote)

	vote := &entity.Vote{
		ID: voteID,
		Voter: &entity.Voter{
			ID: voterID,
		},
		Post: &entity.Post{
			ID: postIDs[0],
		},
		Kind: entity.DownVote,
	}

	err := vs.voteRepository.Update(dbCtx, vote)

	vs.Require().NoError(err)

	updatedVote, err := vs.voteRepository.FindByVoterIDAndPostID(dbCtx, vote.Voter.ID, vote.Post.ID)
	vs.Require().NoError(err)
	vs.Require().NotNil(updatedVote)

	vs.Require().EqualValues(entity.DownVote, updatedVote.Kind)
}

func (vs *PgVoteRepositorySuite) TestFindManyByVoterIDAndPostIDs() {
	ctx := context.Background()
	dbCtx := vs.ctxManager.NewDBContext(ctx)

	authorID := helpers.SetupUser(dbCtx, &vs.Suite)
	voterID := helpers.SetupUser(dbCtx, &vs.Suite)
	sideID := helpers.SetupSide(dbCtx, &vs.Suite)
	postIDs := helpers.SetupPosts(dbCtx, &vs.Suite, authorID, sideID, 10)
	voteIDs := helpers.SetupVotes(dbCtx, &vs.Suite, voterID, postIDs, entity.UpVote)

	votes, err := vs.voteRepository.FindManyByVoterIDAndPostIDs(dbCtx, voterID, postIDs)
	vs.Require().NoError(err)
	vs.Require().Len(votes, len(voteIDs))
}

func (vs *PgVoteRepositorySuite) SetupSuite() {
	vs.PostgresContainerSuite.SetupSuite()
	vs.logger = log.NewTestLogger(vs.T())
	vs.ctxManager = db.NewSqlContextManager(vs.logger, vs.DB)
	vs.voteRepository = repository.NewPgVoteRepository(vs.logger)
}

func TestPgVoteRepository(t *testing.T) {
	suite.Run(t, new(PgVoteRepositorySuite))
}
