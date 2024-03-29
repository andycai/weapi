package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"message": "Page not found.",
	})
}

func NotFoundPage(c *fiber.Ctx) error {
	return c.Render("components/404", fiber.Map{})
}
