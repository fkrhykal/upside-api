package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/stretchr/testify/suite"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresContainerSuite struct {
	suite.Suite
	ctx       context.Context
	container *postgres.PostgresContainer
	DB        *sql.DB
	config    *PostgresDBConfig
}

func (s *PostgresContainerSuite) SetupSuite() {
	s.ctx = context.Background()
	s.config = &PostgresDBConfig{
		Username: "upside",
		Password: "password",
		Database: "test",
		Host:     "localhost",
		Port:     5432,
		SSLMode:  "disable",
		Logger:   log.NewTestLogger(s.T()),
	}
	s.setupContainer()

	pg, err := NewPostgresDB(s.config)
	s.Require().NoError(err)
	s.DB = pg
}

func (s *PostgresContainerSuite) setupContainer() {

	container, err := postgres.Run(
		s.ctx,
		"postgres:17.2-alpine3.21",
		postgres.WithDatabase(s.config.Database),
		postgres.WithUsername(s.config.Username),
		postgres.WithPassword(s.config.Password),
		tc.WithLogger(tc.TestLogger(s.T())),
		postgres.BasicWaitStrategies(),
	)
	s.Require().NoError(err)

	con, err := container.ConnectionString(s.ctx)
	s.Require().NoError(err)

	_, err = fmt.Sscanf(con, "postgres://upside:password@localhost:%d/test?",
		&s.config.Port)
	s.Require().NoError(err)

	s.container = container

}

func (s *PostgresContainerSuite) AfterTest(_, _ string) {
	query := `TRUNCATE TABLE memberships CASCADE; TRUNCATE TABLE sides CASCADE; TRUNCATE TABLE users CASCADE;`
	result, err := s.DB.Exec(query)
	s.T().Log(result)
	s.Require().NoError(err)
}

func (s *PostgresContainerSuite) TearDownSuite() {
	s.container.Terminate(s.ctx)
}
