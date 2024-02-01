package post

import "github.com/gofiber/fiber/v2"

func HandleQueryTags(c *fiber.Ctx, obj any, tableName string) (any, error) {
	return QueryTags()
}
