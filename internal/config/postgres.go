package config

import (
	"strconv"

	"github.com/fkrhykal/upside-api/internal/app"
	"github.com/fkrhykal/upside-api/internal/shared/log"
)

func PostgresConfig(logger log.Logger) *app.PostgresDBConfig {
	portEnv := MustGetEnv("DB_PORT")
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		panic("DB_PORT is not valid")
	}
	return &app.PostgresDBConfig{
		Username: MustGetEnv("DB_USER"),
		Password: MustGetEnv("DB_PASSWORD"),
		Database: MustGetEnv("DB_NAME"),
		Host:     MustGetEnv("DB_HOST"),
		Port:     port,
		SSLMode:  "disable",
		Logger:   logger,
	}
}
