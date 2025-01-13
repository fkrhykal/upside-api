package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/fkrhykal/upside-api/internal/account/dto"
	"github.com/fkrhykal/upside-api/internal/account/repository"
	"github.com/fkrhykal/upside-api/internal/account/service"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/validation"
	"github.com/golang-migrate/migrate/v4"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestSignUp(t *testing.T) {
	ctx := context.Background()

	dbName := "test"
	dbUser := "upside"
	dbPassword := "secret"

	container, err := postgres.Run(
		ctx,
		"postgres:17.2-alpine3.21",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		tc.WithLogger(tc.TestLogger(t)),
		postgres.BasicWaitStrategies(),
	)
	defer container.Terminate(ctx)

	if err != nil {
		t.Fatal(err)
	}

	connection, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(connection)

	pg, err := sql.Open("postgres", fmt.Sprintf("%ssslmode=disable", connection))
	if err != nil {
		t.Fatal(err)
	}

	if err := pg.PingContext(ctx); err != nil {
		t.Fatal(err)
	}

	driver, err := pgMigrate.WithInstance(pg, &pgMigrate.Config{})
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
		dbName,
		driver,
	)

	if err != nil {
		t.Fatal(err)
	}

	if err = m.Up(); err != nil {
		t.Fatal(err)
	}

	ctxManager := db.NewSqlContextManager(pg)
	userRepository := repository.NewPgUserRepository()
	validator := validation.NewValidator()
	authService := service.NewAuthServiceImpl(ctxManager, userRepository, validator)

	req := &dto.SignUpRequest{
		Username: "kontakhaikal",
		Password: "password",
	}

	res, err := authService.SignUp(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)

	dbCtx := ctxManager.NewDBContext(ctx)
	user, err := userRepository.FindByUsername(dbCtx, req.Username)
	if err != nil {
		t.Fatal(err)
	}

	if res.ID.String() != user.ID.String() {
		t.Fatal("user id mismatch")
	}
	t.Log(user)
}
