package app

import (
	"database/sql"

	accountRepositories "github.com/fkrhykal/upside-api/internal/account/repository"
	accountRouters "github.com/fkrhykal/upside-api/internal/account/router"
	accountServices "github.com/fkrhykal/upside-api/internal/account/service"
	"github.com/fkrhykal/upside-api/internal/account/utils"
	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/db"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/fkrhykal/upside-api/internal/shared/validation"
	sideRepositories "github.com/fkrhykal/upside-api/internal/side/repository"
	sideRouters "github.com/fkrhykal/upside-api/internal/side/router"
	sideServices "github.com/fkrhykal/upside-api/internal/side/service"
	"github.com/gofiber/fiber/v2"
)

type BootstrapConfig struct {
	*auth.JwtAuthConfig
	DB        *sql.DB
	Validator validation.Validator
	Logger    log.Logger
	Fiber     *fiber.App
}

func Bootstrap(config *BootstrapConfig) {
	ctxManager := db.NewSqlContextManager(config.Logger, config.DB)
	passwordHasher := utils.NewBcryptPasswordHasher()

	userRepository := accountRepositories.NewPgUserRepository(config.Logger)

	sideRepository := sideRepositories.NewPgSideRepository(config.Logger)
	membershipRepository := sideRepositories.NewPgMembershipRepository(config.Logger)

	authProvider := auth.NewJwtAuthProvider(config.Logger, config.JwtAuthConfig)

	authService := accountServices.NewAuthServiceImpl(
		config.Logger, ctxManager, userRepository, config.Validator, passwordHasher, authProvider)
	userService := accountServices.NewUserServiceImpl(config.Logger, ctxManager, userRepository)

	sideService := sideServices.NewSideServiceImpl(
		config.Logger,
		config.Validator,
		ctxManager,
		sideRepository,
		membershipRepository,
	)

	setupRoutes(
		config.Fiber,
		accountRouters.AuthRouter(config.Logger, authService),
		accountRouters.UserRouter(config.Logger, userService),
		sideRouters.SideRouter(config.Logger, sideService, authProvider),
	)

}

func setupRoutes(app *fiber.App, routerProviders ...func(*fiber.App)) {
	for _, routerProvider := range routerProviders {
		routerProvider(app)
	}
}
