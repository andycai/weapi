package middlewares

import (
	"github.com/andycai/werite/library/authentication"
	"github.com/gofiber/fiber/v2"
)

func AuthorizePage(c *fiber.Ctx) error {
	isAuthenticated, _ := authentication.AuthGet(c)

	if isAuthenticated {
		return c.Next()
	}

	return c.Redirect("/admin/login")
}
