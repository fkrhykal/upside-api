package app

import (
	"database/sql"

	"github.com/fkrhykal/upside-api/internal/account/repository"
	"github.com/fkrhykal/upside-api/internal/account/router"
	"github.com/fkrhykal/upside-api/internal/account/service"
	"github.com/fkrhykal/upside-api/internal/account/utils"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/validation"
	"github.com/gofiber/fiber/v2"
)

type BootstrapConfig struct {
	*service.JwtCredentialConfig
	DB        *sql.DB
	Validator validation.Validator
	Logger    log.Logger
	Fiber     *fiber.App
}

func Bootstrap(config *BootstrapConfig) {
	ctxManager := db.NewSqlContextManager(config.Logger, config.DB)
	passwordHasher := utils.NewBcryptPasswordHasher()

	userRepository := repository.NewPgUserRepository(config.Logger)

	credentialService := service.NewJwtCredentialService(config.Logger, config.JwtCredentialConfig)
	authService := service.NewAuthServiceImpl(
		config.Logger, ctxManager, userRepository, config.Validator, passwordHasher, credentialService)
	userService := service.NewUserServiceImpl(config.Logger, ctxManager, userRepository)

	setupRoutes(
		config.Fiber,
		router.AuthRouter(config.Logger, authService),
		router.UserRouter(config.Logger, userService),
	)
}

func setupRoutes(app *fiber.App, routerProviders ...func(*fiber.App)) {
	for _, routerProvider := range routerProviders {
		routerProvider(app)
	}
}
