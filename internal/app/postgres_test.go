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
	if p.DB == nil {
		p.FailNow("unable to create db")
	}
}

func TestPostgres(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
