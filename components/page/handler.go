package page

import (
	"net/http"
	"strconv"

	"github.com/andycai/weapi/administrator/components/page"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
)

func BeforeRenderPage(ctx *fiber.Ctx, vptr any) (any, error) {
	draft, _ := strconv.ParseBool(ctx.Query("draft"))
	result := vptr.(*model.Page)
	if !draft && !result.Published {
		return nil, core.Error(ctx, http.StatusTooEarly, enum.ErrPageIsNotPublish)
	}
	if draft {
		result.Body = result.Draft
	}
	return page.NewRenderContentFromPage(result), nil
}
