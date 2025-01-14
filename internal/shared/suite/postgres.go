package suite

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/stretchr/testify/suite"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/golang-migrate/migrate/v4"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
)

type PostgresContainerSuite struct {
	suite.Suite
	Ctx       context.Context
	container *postgres.PostgresContainer
	dbName    string
	DB        *sql.DB
}

func (s *PostgresContainerSuite) SetupSuite() {
	s.Ctx = context.Background()
	s.setupContainer()
	s.setupDB()
	s.setupMigration()
}

func (s *PostgresContainerSuite) setupContainer() {
	t := s.T()
	dbName := "test"
	dbUser := "upside"
	dbPassword := "secret"

	container, err := postgres.Run(
		s.Ctx,
		"postgres:17.2-alpine3.21",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		tc.WithLogger(tc.TestLogger(t)),
		postgres.BasicWaitStrategies(),
	)
	s.container = container
	s.dbName = dbName
	if err != nil {
		t.Fatal(err)
	}
}

func (s *PostgresContainerSuite) setupDB() {
	t := s.T()
	connection, err := s.container.ConnectionString(s.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(connection)

	pg, err := sql.Open("postgres", fmt.Sprintf("%ssslmode=disable", connection))
	if err != nil {
		t.Fatal(err)
	}

	if err := pg.PingContext(s.Ctx); err != nil {
		t.Fatal(err)
	}
	s.DB = pg
}

func (s *PostgresContainerSuite) setupMigration() {
	t := s.T()
	driver, err := pgMigrate.WithInstance(s.DB, &pgMigrate.Config{})
	if err != nil {
		t.Fatal(err)
	}

	p, err := filepath.Abs("../../../migrations")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p)

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", filepath.ToSlash(p)),
		s.dbName,
		driver,
	)

	if err != nil {
		t.Fatal(err)
	}

	if err = m.Up(); err != nil {
		t.Fatal(err)
	}
}

func (s *PostgresContainerSuite) TearDownSuite() {
	s.container.Terminate(s.Ctx)
}
