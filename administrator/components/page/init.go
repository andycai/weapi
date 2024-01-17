package page

import (
	"github.com/andycai/weapi"
	"github.com/andycai/weapi/administrator/components/user"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/model"
	"github.com/andycai/weapi/object"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.page.gorm.db"
	KeyNoCheckRouter = "admin.page.router.nocheck"
	KeyCheckRouter   = "admin.page.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initCheckRouter(r fiber.Router) {
	r.Get("/pages/manager", ManagerPage)

	r.Get("/pages/editor", EditorPage)
	r.Get("/pages/editor/:id", EditorPage)
	r.Post("/pages/editor", Create)
	r.Post("/pages/editor/:id", Update)
	r.Get("/pages/movetotrash/:id", MoveToTrashByID)
	r.Post("/pages/movetotrash", MoveToTrash)
	r.Get("/pages/restore/:id", RestoreByID)
	r.Post("/pages/restore", Restore)
}

func initAdminObject() []object.AdminObject {
	return []object.AdminObject{
		{
			Model:       &model.Page{},
			Group:       "Contents",
			Name:        "Page",
			Desc:        "The page data of the website can only be in JSON/YAML format",
			Shows:       []string{"ID", "Site", "Title", "Author", "IsDraft", "Published", "PublishedAt", "CategoryID", "Tags", "CreatedAt"},
			Editables:   []string{"ID", "Site", "CategoryID", "CategoryPath", "Author", "IsDraft", "Draft", "Published", "PublishedAt", "ContentType", "Thumbnail", "Tags", "Title", "Alt", "Description", "Keywords", "Draft", "Remark"},
			Filterables: []string{"Site", "CategoryID", "Tags", "Published", "UpdatedAt"},
			Orderables:  []string{"UpdatedAt", "PublishedAt"},
			Searchables: []string{"ID", "Tags", "Title", "Alt", "Description", "Keywords", "Body"},
			Requireds:   []string{"ID", "Site", "CategoryID", "ContentType", "Body"},
			Icon:        weapi.ReadIcon("./icon/piece.svg"),
			Styles: []string{
				"./css/jsoneditor-9.10.2.min.css",
			},
			Scripts: []object.AdminScript{
				{Src: "./js/cms_widget.js"},
				{Src: "./js/jsoneditor-9.10.2.min.js"},
				{Src: "./js/cms_page.js", Onload: true}},
			Attributes: map[string]object.AdminAttribute{
				"ContentType": {Choices: weapi.EnabledPageContentTypes, Default: enum.ContentTypeJson},
				"Draft":       {Default: "{}"},
				"IsDraft":     {Widget: "is-draft"},
				"Published":   {Widget: "is-published"},
				"Tags":        {Widget: "tags", FilterWidget: "tags"},
				"CategoryID":  {Widget: "category-id-and-path", FilterWidget: "category-id-and-path"},
				"ID":          {Help: "ID must be unique,recommend use page url eg: about-us"},
			},
			EditPage: "./edit_page.html",
			Orders: []object.Order{
				{
					Name: "UpdatedAt",
					Op:   object.OrderOpDesc,
				},
			},
			Actions: []object.AdminAction{
				{
					WithoutObject: true,
					Path:          "save_draft",
					Name:          "Safe Draft",
					Handler: func(db *gorm.DB, c *fiber.Ctx, obj any) (any, error) {
						// return m.handleSaveDraft(db, c, obj)
						return nil, nil
					},
				},
				{
					Path: "duplicate",
					Name: "Duplicate",
					Handler: func(db *gorm.DB, c *fiber.Ctx, obj any) (any, error) {
						// return m.handleMakePageDuplicate(db, c, obj)
						return nil, nil
					},
				},
				{
					Path: "make_publish",
					Name: "Make Publish",
					Handler: func(db *gorm.DB, c *fiber.Ctx, obj any) (any, error) {
						// return m.handleMakePagePublish(db, c, obj, true)
						return nil, nil
					},
				},
				{
					Path: "make_un_publish",
					Name: "Make UnPublish",
					Handler: func(db *gorm.DB, c *fiber.Ctx, obj any) (any, error) {
						// return m.handleMakePagePublish(db, c, obj, false)
						return nil, nil
					},
				},
				{
					WithoutObject: true,
					Path:          "tags",
					Name:          "Query All Tags",
					Handler: func(db *gorm.DB, c *fiber.Ctx, obj any) (any, error) {
						// return m.handleQueryTags(db, c, obj, "pages")
						return nil, nil
					},
				},
			},
			BeforeCreate: func(db *gorm.DB, ctx *fiber.Ctx, vptr any) error {
				page := vptr.(*model.Page)
				page.ContentType = enum.ContentTypeJson
				page.Creator = *user.CurrentUser(ctx)
				page.IsDraft = true
				return nil
			},
			BeforeUpdate: func(db *gorm.DB, ctx *fiber.Ctx, vptr any, vals map[string]any) error {
				page := vptr.(*model.Page)
				page.IsDraft = true
				if _, ok := vals["published"]; ok {
					page.Published = vals["published"].(bool)
					if page.Published {
						page.Body = page.Draft
						page.IsDraft = false
					}
				}
				return nil
			},
		},
	}
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initCheckRouter)
	core.RegisterAdminObject(initAdminObject())
}
