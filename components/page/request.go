package page

import (
	"database/sql"
	"time"

	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
)

type requestCreate struct {
	ID          uint   `json:"id"`
	Slug        string `json:"slug"`
	Title       string `json:"title" validate:"required"`
	Body        string `json:"body" validate:"required"`
	PublishedAt string `json:"published_at" form:"published_at" validate:"required"`
}

func Bind(c *fiber.Ctx, page *model.Page) error {
	var r requestCreate
	if err := c.BodyParser(&r); err != nil {
		return err
	}

	if err := core.Validate(r); err != nil {
		return err
	}

	// page.ID = r.ID
	page.Title = r.Title
	page.Body = r.Body

	if r.Slug != "" {
		page.ID = r.Slug
	} else {
		page.ID = slug.Make(r.Title)
	}

	if r.PublishedAt != "" {
		page.PublishedAt = sql.NullTime{Time: core.ParseDate(r.PublishedAt), Valid: true} // core.ParseDate(r.PublishedAt)
	} else {
		page.PublishedAt = sql.NullTime{Time: time.Now(), Valid: true} // time.Now()
	}

	return nil
}
