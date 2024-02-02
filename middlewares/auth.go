package middlewares

import (
	"github.com/andycai/weapi/lib/authentication"
	"github.com/gofiber/fiber/v2"
)

func Authorize(c *fiber.Ctx) error {
	isAuthenticated, _ := authentication.AuthGet(c)

	if isAuthenticated {
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": true,
		"msg":   "Unauthorized",
	})
}

func AuthorizePage(c *fiber.Ctx) error {
	isAuthenticated, _ := authentication.AuthGet(c)

	if isAuthenticated {
		return c.Next()
	}

	return c.Redirect("/auth/login")
}
