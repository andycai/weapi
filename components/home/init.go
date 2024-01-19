package home

import (
	"github.com/andycai/weapi/core"
	"github.com/gofiber/fiber/v2"
)

const (
	KeyHomeNoCheckRouter = "home.router.nocheck"
)

func initNoCheckRouter(r fiber.Router) {
	r.Get("/", HomePage)

	// HTMX
	r.Get("/htmx/home", HTMXHomePage)
	r.Get("/htmx/home/tag-list", HTMXHomeTagList)
	r.Get("/htmx/home/global-feed", HTMXHomeGlobalFeed)
}

func init() {
	core.RegisterNoCheckRouter(KeyHomeNoCheckRouter, initNoCheckRouter)
}
