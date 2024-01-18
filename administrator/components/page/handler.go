package page

import (
	"fmt"
	"html/template"

	"github.com/andycai/weapi/administrator/utils"
	"github.com/andycai/weapi/components/page"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/library/authentication"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
)

func ManagerPage(c *fiber.Ctx) error {
	var (
		total             int64
		totalAll          int64
		totalTrash        int64
		voList            []model.Page
		q                 string
		status            string
		queryParam        string
		CurrentPagination int
	)
	q = c.Query("q")
	status = c.Query("status")

	if c.QueryInt("page") > 1 {
		CurrentPagination = c.QueryInt("page") - 1
	}

	totalAll = page.Dao.Count()
	totalTrash = page.Dao.CountTrash()

	switch status {
	case "trash":
		voList = page.Dao.GetTrashListByPage(CurrentPagination, 20, q)
		total = totalTrash
	default:
		voList = page.Dao.GetListByPage(CurrentPagination, 20, q)
		total = totalAll
	}

	if status != "" {
		queryParam = fmt.Sprintf("&status=%s", status)
	}

	totalPagination, hasPagination := utils.CalcPagination(total)

	return core.Render(c, "admin/pages/pages", fiber.Map{
		"PageTitle":         "All Pages",
		"NavBarActive":      "pages",
		"Path":              "/admin/pages/manager",
		"Pages":             voList,
		"Q":                 q,
		"Status":            status,
		"Total":             total,
		"TotalAll":          totalAll,
		"TotalTrash":        totalTrash,
		"TotalPagination":   totalPagination,
		"HasPagination":     hasPagination,
		"CurrentPagination": CurrentPagination + 1,
		"QueryParam":        template.URL(queryParam),
	}, "admin/layouts/app")
}

func EditorPage(c *fiber.Ctx) error {
	var pageVo model.Page
	hasPage := false

	if c.Params("id") != "" {
		id := cast.ToUint(c.Params("id"))
		hasPage = true
		vo, _ := page.Dao.GetByID(id)
		pageVo = *vo
	}

	return core.Render(c, "admin/pages/page", fiber.Map{
		"PageTitle":    "Page Editor",
		"NavBarActive": "pages",
		"Path":         "/admin/pages/editor",
		"Domain":       "127.0.0.1",
		"HasPage":      hasPage,
		"Page":         pageVo,
	}, "admin/layouts/app")
}

func Create(c *fiber.Ctx) error {
	var pageVo model.Page

	err := page.Bind(c, &pageVo)
	if err != nil {
		return err
	}

	_, userID := authentication.AuthGet(c)
	pageVo.CreatorID = userID

	db.Create(&pageVo)

	core.PushMessages(fmt.Sprintf("Created page id:%d, title:%s", pageVo.ID, pageVo.Title))

	return c.Redirect("/admin/pages/manager")
}

func Update(c *fiber.Ctx) error {
	var pageVo model.Page

	err := page.Bind(c, &pageVo)
	if err != nil {
		return err
	}

	db.Omit("created_at", "user_id").Save(&pageVo)

	core.PushMessages(fmt.Sprintf("Updated page id:%d, title:%s", pageVo.ID, pageVo.Title))

	return c.Redirect("/admin/pages/manager")
}

func MoveToTrashByID(c *fiber.Ctx) error {
	id := cast.ToUint(c.Params("id"))
	if id > 0 {
		page.Dao.DeleteByIds([]uint{id})
		core.PushMessages(fmt.Sprintf("Move to trash: %v", id))
	}

	return c.Redirect("/admin/pages/manager")
}

func MoveToTrash(c *fiber.Ctx) error {
	form := &utils.FormIDArray{}
	if err := c.BodyParser(form); err != nil {
		return err
	}
	if len(form.ID) > 0 {
		page.Dao.DeleteByIds(form.ID)
		core.PushMessages(fmt.Sprintf("Move to trash: %v", form.ID))
	}

	return c.Redirect("/admin/pages/manager")
}

func DeletePermanetly(c *fiber.Ctx) error {
	form := &utils.FormIDArray{}
	if err := c.BodyParser(form); err != nil {
		return err
	}
	if len(form.ID) > 0 {
		page.Dao.DeletePermanetlyByIds(form.ID)
		core.PushMessages(fmt.Sprintf("Delete permanently: %v", form.ID))
	}

	return c.Redirect("/admin/pages/manager")
}

func RestoreByID(c *fiber.Ctx) error {
	id := cast.ToUint(c.Params("id"))
	if id > 0 {
		page.Dao.RestoreByIds([]uint{id})
		core.PushMessages(fmt.Sprintf("Restore pages: %v", id))
	}

	return c.Redirect("/admin/pages/manager")
}

func Restore(c *fiber.Ctx) error {
	form := &utils.FormIDArray{}
	if err := c.BodyParser(form); err != nil {
		return err
	}
	if len(form.ID) > 0 {
		page.Dao.RestoreByIds(form.ID)
		core.PushMessages(fmt.Sprintf("Restore pages: %v", form.ID))
	}

	return c.Redirect("/admin/pages/manager")
}

//#region action handler

func HandleMakePagePublish(c *fiber.Ctx, obj any, publish bool) (any, error) {
	siteId := c.Query("site_id")
	id := c.Query("id")
	if err := model.MakePublish(db, siteId, id, obj, publish); err != nil {
		// carrot.Warning("make publish failed:", siteId, id, publish, err)
		return false, err
	}
	return true, nil
}

func HandleMakePageDuplicate(c *fiber.Ctx, obj any) (any, error) {
	if err := model.MakeDuplicate(db, obj); err != nil {
		// carrot.Warning("make duplicate failed:", obj, err)
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

	if err := model.SafeDraft(db, siteId, id, obj, draft); err != nil {
		// carrot.Warning("safe draft failed:", siteId, id, err)
		return false, err
	}
	return true, nil
}

func HandleQueryTags(c *fiber.Ctx, obj any, tableName string) (any, error) {
	return model.QueryTags(db.Table(tableName))
}

//#endregion
