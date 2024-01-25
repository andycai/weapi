package post

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
	Description string `json:"description" validate:"required"`
	Body        string `json:"body" validate:"required"`
	Action      string `json:"action" form:"action"`
	CategoryID  string `json:"category_id" form:"category_id"`
	PublishedAt string `json:"published_at" form:"published_at" validate:"required"`
}

func Bind(c *fiber.Ctx, post *model.Post) error {
	var r requestCreate
	if err := c.BodyParser(&r); err != nil {
		return err
	}

	if err := core.Validate(r); err != nil {
		return err
	}

	// post.ID = r.ID
	post.Title = r.Title
	post.Description = r.Description
	post.Body = r.Body
	post.CategoryID = r.CategoryID

	if r.Action == "draft" {
		post.IsDraft = true
	} else {
		post.IsDraft = false
	}

	if r.Slug != "" {
		post.ID = r.Slug
	} else {
		post.ID = slug.Make(r.Title)
	}

	if r.PublishedAt != "" {
		post.PublishedAt = sql.NullTime{Time: core.ParseDate(r.PublishedAt), Valid: true} // core.ParseDate(r.PublishedAt)
	} else {
		post.PublishedAt = sql.NullTime{Time: time.Now(), Valid: true} // time.Now()
	}

	return nil
}

type requestCategoryCreate struct {
	ID          uint   `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func BindCategory(c *fiber.Ctx, category *model.Category) error {
	var r requestCategoryCreate
	if err := c.BodyParser(&r); err != nil {
		return err
	}

	if err := core.Validate(r); err != nil {
		return err
	}

	// category.ID = r.ID
	category.Name = r.Name
	// category.Description = r.Description

	if r.Slug != "" {
		category.UUID = r.Slug
	} else {
		category.UUID = slug.Make(r.Name)
	}

	return nil
}
