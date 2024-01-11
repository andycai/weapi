package post

import (
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
	CategoryID  uint   `json:"category_id" form:"category_id"`
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

	post.ID = r.ID
	post.Title = r.Title
	post.Description = r.Description
	post.Body = r.Body
	post.CategoryID = r.CategoryID

	if r.Action == "draft" {
		post.IsDraft = 1
	} else {
		post.IsDraft = 0
	}

	if r.Slug != "" {
		post.Slug = r.Slug
	} else {
		post.Slug = slug.Make(r.Title)
	}

	if r.PublishedAt != "" {
		post.PublishedAt = core.ParseDate(r.PublishedAt)
	} else {
		post.PublishedAt = time.Now()
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

	category.ID = r.ID
	category.Name = r.Name
	category.Description = r.Description

	if r.Slug != "" {
		category.Slug = r.Slug
	} else {
		category.Slug = slug.Make(r.Name)
	}

	return nil
}

type requestTagCreate struct {
	ID   uint   `json:"id"`
	Name string `json:"name" validate:"required"`
	Slug string `json:"slug"`
}

func BindTag(c *fiber.Ctx, tag *model.Tag) error {
	var r requestTagCreate
	if err := c.BodyParser(&r); err != nil {
		return err
	}

	if err := core.Validate(r); err != nil {
		return err
	}

	tag.ID = r.ID
	tag.Name = r.Name

	if r.Slug != "" {
		tag.Slug = r.Slug
	} else {
		tag.Slug = slug.Make(r.Name)
	}

	return nil
}
