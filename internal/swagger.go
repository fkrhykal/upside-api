package internal

import (
	"log"

	"github.com/fkrhykal/upside-api/internal/account/router"
)

func main() {
	// @title Upside API
	// @version 1.0
	// @description This is api for upside application
	// @BasePath /
	generateSwagger(
		router.AuthRouter,
	)
}
func generateSwagger(routers ...any) {
	for _, router := range routers {
		log.Print(router)
	}
}
