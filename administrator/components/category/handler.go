package category

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func HandleQueryCategoryWithCount(c *fiber.Ctx, obj any) (any, error) {
	siteId := c.Query("site_id")
	current := strings.ToLower(c.Query("current"))
	return QueryCategoryWithCount(siteId, current)
}
