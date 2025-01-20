package main

import (
	"log"

	_ "github.com/fkrhykal/upside-api/docs/api"
	"github.com/fkrhykal/upside-api/internal/app"
	"github.com/fkrhykal/upside-api/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}
	logger := app.NewLogrus()
	fiber := app.NewFiber(logger)
	validator := app.NewGoPlaygroundValidator(logger)
	pg, err := app.NewPostgresDB(config.PostgresConfig(logger))
	if err != nil {
		log.Fatal(err)
	}

	app.Bootstrap(&app.BootstrapConfig{
		Fiber:               fiber,
		Logger:              logger,
		DB:                  pg,
		Validator:           validator,
		JwtCredentialConfig: config.DefaultJwtCredentialConfig(),
	})

	if err := fiber.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
