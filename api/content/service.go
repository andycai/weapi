package content

import (
	"net/http"
	"strconv"

	"github.com/andycai/weapi/admin/content"
	"github.com/andycai/weapi/admin/user"
	"github.com/andycai/weapi/constant"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
)

func BeforeRenderPage(ctx *fiber.Ctx, vptr any) (any, error) {
	draft, _ := strconv.ParseBool(ctx.Query("draft"))
	result := vptr.(*model.Page)
	if !draft && !result.Published {
		return nil, core.Error(ctx, http.StatusTooEarly, constant.ErrPageIsNotPublish)
	}
	if draft {
		result.Body = result.Draft
	}
	return content.NewRenderContentFromPage(result), nil
}

func BeforeRenderPost(ctx *fiber.Ctx, vptr any) (any, error) {
	draft, _ := strconv.ParseBool(ctx.Query("draft"))
	result := vptr.(*model.Post)
	if !draft && !result.Published {
		return nil, constant.ErrPostIsNotPublish
	}
	if draft {
		result.Body = result.Draft
	}

	relations := true
	if ctx.Method() == http.MethodPost { // batch query
		relations = false
	}

	return content.NewRenderContentFromPost(result, relations), nil
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

	relationCount := user.GetIntValue(constant.KEY_CMS_RELATION_COUNT, 3)
	suggestionCount := user.GetIntValue(constant.KEY_CMS_SUGGESTION_COUNT, 3)

	r.Suggestions, _ = content.GetSuggestions(siteId, categoryId, categoryPath, "", relationCount)
	r.Relations, _ = content.GetRelations(siteId, categoryId, categoryPath, "", suggestionCount)

	return r, nil
}
