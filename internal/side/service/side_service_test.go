package service_test

import (
	"context"
	"testing"

	"github.com/fkrhykal/upside-api/internal/app"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/log"
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
		FounderID:   userID,
	}
	_, err = s.CreateSide(ctx, req)
	s.Nil(err)
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
