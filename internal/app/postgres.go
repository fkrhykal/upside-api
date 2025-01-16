package app

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

type PostgresDB = sql.DB

type PostgresDBConfig struct {
	Username string
	Password string
	Database string
	Host     string
	Port     int
	SSLMode  string
	Logger   log.Logger
}

func NewPostgresDB(config *PostgresDBConfig) (*PostgresDB, error) {
	config.Logger.Debugf("Initializing PostgreSQL connection with config: %+v", config)

	dataSource := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%+v",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.SSLMode,
	)

	pg, err := sql.Open("postgres", dataSource)
	if err != nil {
		config.Logger.Errorf("Failed to open PostgreSQL connection: %+v", err)
		return nil, err
	}
	config.Logger.Info("PostgreSQL connection opened successfully")

	if err := pg.Ping(); err != nil {
		config.Logger.Errorf("Failed to ping PostgreSQL database: %+v", err)
		return nil, err
	}
	config.Logger.Info("PostgreSQL database connection verified")

	config.Logger.Info("Setting up PostgreSQL migrations")
	if err := setupPostgresMigration(&migrationConfig{
		DB:     pg,
		DBName: config.Database,
		Logger: config.Logger,
	}); err != nil {
		config.Logger.Errorf("Failed to set up PostgreSQL migrations: %+v", err)
		pg.Close()
		return nil, err
	}
	config.Logger.Info("PostgreSQL migrations completed successfully")

	return pg, nil
}

type migrationConfig struct {
	DB     *PostgresDB
	DBName string
	Logger log.Logger
}

func setupPostgresMigration(config *migrationConfig) error {
	config.Logger.Debugf("Starting migration setup for database: %s", config.DBName)

	driver, err := postgres.WithInstance(config.DB, &postgres.Config{})
	if err != nil {
		config.Logger.Errorf("Failed to create database driver: %+v", err)
		return err
	}
	config.Logger.Debug("Database driver created successfully")

	sourceDriver, err := iofs.New(migrations.MigrationFs, ".")
	if err != nil {
		config.Logger.Errorf("Failed to create source driver for migrations: %+v", err)
		return err
	}
	config.Logger.Debug("Source driver for migrations created successfully")

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
	if err != nil {
		config.Logger.Errorf("Failed to initialize migration instance: %+v", err)
		return err
	}
	config.Logger.Debug("Migration instance initialized successfully")

	config.Logger.Info("Applying database migrations")

	err = m.Up()

	if err == nil {
		config.Logger.Info("Migrations applied successfully")
		return nil
	}

	if errors.Is(err, migrate.ErrNoChange) {
		config.Logger.Info("No new migrations to apply")
		return nil
	}

	config.Logger.Errorf("Migration process failed: %+v", err)
	return err
}
