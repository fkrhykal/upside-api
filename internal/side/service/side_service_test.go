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

	sides, err := s.SideService.GetJoinedSides(authCtx, &pagination.OffsetBased{Page: 1, Limit: 2})
	s.Nil(err)
	s.Len(sides, 2)

	offsetSides, err := s.SideService.GetJoinedSides(authCtx, &pagination.OffsetBased{Page: 2, Limit: 2})
	s.Nil(err)
	s.Len(offsetSides, 2)

	for i := range 2 {
		s.NotNil(sides[i].MembershipDetail)
		s.NotNil(offsetSides[i].MembershipDetail)
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

	sides, err := s.SideService.GetPopularSides(authCtx, &pagination.OffsetBased{Page: 1, Limit: 2})
	s.Nil(err)
	s.Len(sides, 2)

	offsetSides, err := s.SideService.GetPopularSides(authCtx, &pagination.OffsetBased{Page: 2, Limit: 2})
	s.Nil(err)
	s.Len(offsetSides, 2)

	for i := range 2 {
		s.NotNil(sides[i].MembershipDetail)
		s.NotNil(offsetSides[i].MembershipDetail)
	}
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

	sides, err := s.SideService.GetSides(authCtx, &pagination.OffsetBased{Page: 1, Limit: 2})
	s.Nil(err)
	s.Len(sides, 2)

	offsetSides, err := s.SideService.GetSides(authCtx, &pagination.OffsetBased{Page: 2, Limit: 2})
	s.Nil(err)
	s.Len(offsetSides, 2)

	for i := range 2 {
		s.NotNil(sides[i].MembershipDetail)
		s.NotNil(offsetSides[i].MembershipDetail)
	}
}

func (s *SideServiceSuite[T]) saveUser() uuid.UUID {
	userID := uuid.New()
	query := `INSERT INTO users(id, username, password) VALUES($1, $2, $3)`
	_, err := s.DB.Exec(query, userID, faker.Username(), faker.Password())
	s.Nil(err)
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
