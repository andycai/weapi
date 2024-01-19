package post

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/andycai/weapi/components/user"
	"github.com/andycai/weapi/conf"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/library/authentication"
	"github.com/andycai/weapi/model"
	"github.com/andycai/weapi/object"
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

func GetPostOrPageDB(ctx *fiber.Ctx, isCreate bool) *gorm.DB {
	if isCreate {
		return db
	}
	draft, _ := strconv.ParseBool(ctx.Query("draft"))
	if draft {
		return db
	}
	// single get not need published
	if ctx.Method() == http.MethodGet {
		return db
	}
	// query must be published
	return db.Where("published", true)
}

func BeforeRenderPost(ctx *fiber.Ctx, vptr any) (any, error) {
	draft, _ := strconv.ParseBool(ctx.Query("draft"))
	result := vptr.(*model.Post)
	if !draft && !result.Published {
		return nil, enum.ErrPostIsNotPublish
	}
	if draft {
		result.Body = result.Draft
	}

	relations := true
	if ctx.Method() == http.MethodPost { // batch query
		relations = false
	}

	return model.NewRenderContentFromPost(db, result, relations), nil
}

func BeforeQueryRenderPost(ctx *fiber.Ctx, queryResult *object.QueryResult) (any, error) {
	if len(queryResult.Items) <= 0 {
		return nil, nil
	}
	firstItem, ok := queryResult.Items[0].(*model.RenderContent)
	if !ok {
		return nil, nil
	}

	siteId := firstItem.SiteID
	categoryId := ""
	categoryPath := ""
	if firstItem.Category != nil {
		categoryId = firstItem.Category.UUID
		categoryPath = firstItem.Category.Path
	}

	r := &model.ContentQueryResult{
		QueryResult: queryResult,
	}

	relationCount := conf.GetIntValue(db, enum.KEY_CMS_RELATION_COUNT, 3)
	suggestionCount := conf.GetIntValue(db, enum.KEY_CMS_SUGGESTION_COUNT, 3)

	r.Suggestions, _ = model.GetSuggestions(db, siteId, categoryId, categoryPath, "", relationCount)
	r.Relations, _ = model.GetRelations(db, siteId, categoryId, categoryPath, "", suggestionCount)

	return r, nil
}
