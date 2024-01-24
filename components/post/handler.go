package post

import (
	"net/http"
	"strconv"

	"github.com/andycai/weapi/conf"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/model"
	"github.com/andycai/weapi/object"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

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
