package page

import (
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/log"
	"github.com/gofiber/fiber/v2"
)

//#region action handler

func HandleMakePagePublish(c *fiber.Ctx, obj any, publish bool) (any, error) {
	siteId := c.Query("site_id")
	id := c.Query("id")
	if err := MakePublish(siteId, id, obj, publish); err != nil {
		log.Infof("make publish failed: %s, %s, %t, %s", siteId, id, publish, err)
		return false, err
	}
	return true, nil
}

func HandleMakePageDuplicate(c *fiber.Ctx, obj any) (any, error) {
	if err := MakeDuplicate(obj); err != nil {
		log.Infof("make duplicate failed: %v, %s", obj, err)
		return false, err
	}
	return true, nil
}

func HandleSaveDraft(c *fiber.Ctx, obj any) (any, error) {
	siteId := c.Query("site_id")
	id := c.Query("id")

	var formData map[string]string
	if err := c.BodyParser(&formData); err != nil {
		return nil, err
	}

	draft, ok := formData["draft"]
	if !ok {
		return nil, enum.ErrDraftIsInvalid
	}

	if err := SafeDraft(siteId, id, obj, draft); err != nil {
		log.Infof("safe draft failed: %s, %s, %s", siteId, id, err)
		return false, err
	}
	return true, nil
}

func HandleQueryTags(c *fiber.Ctx, obj any, tableName string) (any, error) {
	return QueryTags()
}

//#endregion
