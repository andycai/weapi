package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func NotFoundPage(c *fiber.Ctx) error {
	return c.Render("components/404", fiber.Map{})
}
