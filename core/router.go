package core

import (
	"github.com/andycai/weapi/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {
	page := app.Group("/")
	for _, f := range routerRootNoCheckMap {
		f(page)
	}

	api := app.Group("/api", middlewares.Authorize)
	for _, f := range routerAPICheckMap {
		f(api)
	}

	admin := app.Group("/admin", middlewares.Authorize)
	for _, f := range routerAdminCheckMap {
		f(admin)
	}

	app.Use(middlewares.NotFound)
}
