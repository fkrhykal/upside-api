package app_test

import (
	"testing"

	"github.com/fkrhykal/upside-api/internal/app"
	"github.com/stretchr/testify/suite"
)

type PostgresTestSuite struct {
	app.PostgresContainerSuite
}

func (p *PostgresTestSuite) TestCreatePostgresDB() {
	p.Require().NotNil(p.DB, "unable to create DB")
}

func TestPostgres(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
