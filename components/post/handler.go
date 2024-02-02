package post

import (
	"net/http"
	"strconv"

	"github.com/andycai/weapi/administrator/components/config"
	"github.com/andycai/weapi/administrator/components/post"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
)

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

	return post.NewRenderContentFromPost(result, relations), nil
}

func BeforeQueryRenderPost(ctx *fiber.Ctx, queryResult *model.QueryResult) (any, error) {
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

	relationCount := config.GetIntValue(enum.KEY_CMS_RELATION_COUNT, 3)
	suggestionCount := config.GetIntValue(enum.KEY_CMS_SUGGESTION_COUNT, 3)

	r.Suggestions, _ = post.GetSuggestions(siteId, categoryId, categoryPath, "", relationCount)
	r.Relations, _ = post.GetRelations(siteId, categoryId, categoryPath, "", suggestionCount)

	return r, nil
}
