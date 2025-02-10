package service_test

import (
	"context"
	"testing"

	"github.com/fkrhykal/upside-api/internal/app"
	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/pagination"
	"github.com/fkrhykal/upside-api/internal/side/dto"
	"github.com/fkrhykal/upside-api/internal/side/repository"
	"github.com/fkrhykal/upside-api/internal/side/service"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type SideServiceSuite[T any] struct {
	app.PostgresContainerSuite
	service.SideService
}

func (s *SideServiceSuite[T]) TestCreateSide() {
	ctx := context.Background()

	userID := uuid.New()

	query := "INSERT INTO users(id, username, password) VALUES($1, $2, $3)"
	_, err := s.DB.Exec(query, userID, faker.Username(), faker.Password())
	s.Nil(err)

	req := &dto.CreateSideRequest{
		Nick:        faker.Username(),
		Name:        faker.Name(),
		Description: faker.Sentence(),
	}

	authCtx := auth.NewAuthContext(ctx, &auth.UserCredential{ID: userID})
	_, err = s.CreateSide(authCtx, req)
	s.Nil(err)
}

func (s *SideServiceSuite[T]) TestGetJoinedSides() {
	userID := s.saveUser()
	ctx := context.Background()
	authCtx := auth.NewAuthContext(ctx, &auth.UserCredential{ID: userID})

	for range 100 {
		req := &dto.CreateSideRequest{
			Nick:        faker.Username(),
			Name:        faker.Name(),
			Description: faker.Sentence(),
		}

		_, err := s.CreateSide(authCtx, req)
		s.Nil(err)
	}

	res, err := s.SideService.GetJoinedSides(authCtx, &pagination.OffsetBased{Page: 1, Limit: 2})
	s.Nil(err)
	s.Len(res.Sides, 2)

	offsetRes, err := s.SideService.GetJoinedSides(authCtx, &pagination.OffsetBased{Page: 2, Limit: 2})
	s.Nil(err)
	s.Len(offsetRes.Sides, 2)

	for i := range 2 {
		s.NotNil(res.Sides[i].MembershipDetail)
		s.NotNil(offsetRes.Sides[i].MembershipDetail)
		s.EqualValues(50, res.Metadata.TotalPage)
		s.EqualValues(50, offsetRes.Metadata.TotalPage)
	}
}

func (s *SideServiceSuite[T]) TestGetPopularSides() {
	userID := s.saveUser()
	ctx := context.Background()
	authCtx := auth.NewAuthContext(ctx, &auth.UserCredential{ID: userID})

	for range 100 {
		req := &dto.CreateSideRequest{
			Nick:        faker.Username(),
			Name:        faker.Name(),
			Description: faker.Sentence(),
		}

		_, err := s.CreateSide(authCtx, req)
		s.Nil(err)
	}

	res, err := s.SideService.GetPopularSides(authCtx, &pagination.OffsetBased{Page: 1, Limit: 2})
	s.Nil(err)
	s.Len(res.Sides, 2)

	offsetRes, err := s.SideService.GetPopularSides(authCtx, &pagination.OffsetBased{Page: 2, Limit: 2})
	s.Nil(err)
	s.Len(offsetRes.Sides, 2)

	s.EqualValues(50, res.Metadata.TotalPage)
	s.EqualValues(50, offsetRes.Metadata.TotalPage)
}

func (s *SideServiceSuite[T]) TestGetSides() {
	userID := s.saveUser()
	ctx := context.Background()
	authCtx := auth.NewAuthContext(ctx, &auth.UserCredential{ID: userID})

	for range 100 {
		req := &dto.CreateSideRequest{
			Nick:        faker.Username(),
			Name:        faker.Name(),
			Description: faker.Sentence(),
		}

		_, err := s.CreateSide(authCtx, req)
		s.Nil(err)
	}

	res, err := s.SideService.GetSides(authCtx, &pagination.OffsetBased{Page: 1, Limit: 2})
	s.Nil(err)
	s.Len(res.Sides, 2)

	offsetRes, err := s.SideService.GetSides(authCtx, &pagination.OffsetBased{Page: 2, Limit: 2})
	s.Nil(err)
	s.Len(offsetRes.Sides, 2)

	s.Require().EqualValues(50, res.Metadata.TotalPage)
	s.Require().EqualValues(50, offsetRes.Metadata.TotalPage)

}

func (s *SideServiceSuite[T]) TestJoinSide() {
	founderID := s.saveUser()
	ctx := context.Background()
	founderAuthCtx := auth.NewAuthContext(ctx, &auth.UserCredential{ID: founderID})

	createSideRequest := &dto.CreateSideRequest{
		Nick:        faker.Username(),
		Name:        faker.Name(),
		Description: faker.Sentence(),
	}

	createSideResponse, err := s.CreateSide(founderAuthCtx, createSideRequest)
	s.Require().NoError(err)

	memberID := s.saveUser()
	memberAuthCtx := auth.NewAuthContext(ctx, &auth.UserCredential{ID: memberID})

	joinSideResponse, err := s.JoinSide(memberAuthCtx, &dto.JoinSideRequest{
		SideID: createSideResponse.ID,
	})

	s.Require().NoError(err)
	s.Require().NotNil(joinSideResponse)
	s.Require().EqualValues(createSideResponse.ID, joinSideResponse.SideID)
}

func (s *SideServiceSuite[T]) saveUser() uuid.UUID {
	userID := uuid.New()
	query := `INSERT INTO users(id, username, password) VALUES($1, $2, $3)`
	_, err := s.DB.Exec(query, userID, faker.Username(), faker.Password())
	s.Require().NoError(err)
	return userID
}

func (s *SideServiceSuite[T]) SetupSuite() {
	s.PostgresContainerSuite.SetupSuite()
	logger := log.NewTestLogger(s.T())
	validator := app.NewGoPlaygroundValidator(logger)
	ctxManager := db.NewSqlContextManager(logger, s.DB)
	sideRepository := repository.NewPgSideRepository(logger)
	membershipRepository := repository.NewPgMembershipRepository(logger)
	s.SideService = service.NewSideServiceImpl(
		logger,
		validator,
		ctxManager,
		sideRepository,
		membershipRepository,
	)
}

func TestSideService(t *testing.T) {
	suite.Run(t, new(SideServiceSuite[any]))
}
