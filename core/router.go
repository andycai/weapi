package core

import (
	"github.com/andycai/werite/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {
	page := app.Group("/")
	for _, f := range routerNoCheckMap {
		f(page)
	}

	api := app.Group("/api")
	for _, f := range routerCheckMap {
		f(api)
	}

	admin := app.Group("/admin", middlewares.AuthorizePage)
	for _, f := range routerAdminCheckMap {
		f(admin)
	}

	app.Use(middlewares.NotFound)
}
