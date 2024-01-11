package post

import (
	"errors"

	"github.com/andycai/weapi/components/user"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/library/authentication"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PostDetailPage(c *fiber.Ctx) error {
	var post *model.Post
	var authenticatedUser model.User
	isSelf := false

	isAuthenticated, userID := authentication.AuthGet(c)

	post, err := Dao.GetBySlug(c.Params("slug"))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Redirect("/")
		}
	}

	if isAuthenticated {
		authenticatedUser = *user.Dao.GetByID(userID)
	}

	return core.Render(c, "posts/show", fiber.Map{
		"PageTitle":         post.Title + " — Werite",
		"Post":              post,
		"FiberCtx":          c,
		"IsOob":             false,
		"IsSelf":            isSelf,
		"IsPostFavorited":   false,
		"AuthenticatedUser": authenticatedUser,
	}, "layouts/app")
}

//#region HTMX interface

// HTMXHomePostDetailPage detail page
func HTMXHomePostDetailPage(c *fiber.Ctx) error {
	var post *model.Post
	isSelf := false
	var authenticatedUser *model.User

	isAuthenticated, userID := authentication.AuthGet(c)

	if isAuthenticated {
		authenticatedUser = user.Dao.GetByID(userID)
	}

	post, err := Dao.GetBySlug(c.Params("slug"))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Redirect("/")
		}
	}

	return core.Render(c, "posts/htmx-post-page", fiber.Map{
		"PageTitle":         post.Title,
		"NavBarActive":      "none",
		"Post":              post,
		"IsOob":             false,
		"IsSelf":            isSelf,
		"IsPostFavorited":   false,
		"AuthenticatedUser": authenticatedUser,
		"FiberCtx":          c,
	}, "layouts/app-htmx")
}

//#endregion
