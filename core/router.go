package core

import (
	"github.com/andycai/weapi/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {
	public := app.Group("/")
	for _, f := range routerPublicNoCheckMap {
		f(public)
	}

	api := app.Group("/api", withAPIAuth)
	for _, f := range routerAPICheckMap {
		f(api)
	}

	admin := app.Group("/admin", withAdminAuth)
	for _, f := range routerAdminCheckMap {
		f(admin)
	}

	app.Use(middlewares.NotFound)
}
