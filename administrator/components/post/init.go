package post

import (
	"github.com/andycai/weapi"
	"github.com/andycai/weapi/administrator/components/page"
	"github.com/andycai/weapi/administrator/components/user"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/model"
	"github.com/andycai/weapi/object"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.post.gorm.db"
	KeyNoCheckRouter = "admin.post.router.nocheck"
	KeyCheckRouter   = "admin.post.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initCheckRouter(r fiber.Router) {
	// r.Get("/posts/manager", ManagerPage)
	// r.Get("/posts/editor", EditorPage)
	// r.Get("/posts/editor/:id", EditorPage)
	// r.Post("/posts/editor", Create)
	// r.Post("/posts/editor/:id", Update)
	// r.Get("/posts/movetotrash/:id", MoveToTrashByID)
	// r.Post("/posts/movetotrash", MoveToTrash)
	// r.Get("/posts/restore/:id", RestoreByID)
	// r.Post("/posts/restore", Restore)

	// r.Get("/categories/manager", ManagerCategoryPage)
	// r.Get("/categories/editor", EditorCategoryPage)
	// r.Get("/categories/editor/:id", EditorCategoryPage)
	// r.Post("/categories/editor", CreateCategory)
	// r.Post("/categories/editor/:id", UpdateCategory)
	// r.Post("/categories/delete", DeleteCategories)

	// r.Get("/tags/manager", ManagerTagsPage)
	// r.Get("/tags/editor", EditorTagPage)
	// r.Post("/tags/editor", CreateTag)
	// r.Post("/tags/delete", DeleteTags)
}

func initAdminObject() []object.AdminObject {
	return []object.AdminObject{
		{
			Model:       &model.Post{},
			Group:       "Contents",
			Name:        "Post",
			PluralName:  "Posts",
			Desc:        "Website articles or blogs, support HTML and Markdown formats",
			Shows:       []string{"ID", "Site", "Title", "Author", "CategoryID", "Tags", "IsDraft", "Published", "PublishedAt", "CreatedAt"},
			Editables:   []string{"ID", "Site", "CategoryID", "CategoryPath", "Author", "IsDraft", "Draft", "Published", "PublishedAt", "ContentType", "Thumbnail", "Tags", "Title", "Alt", "Description", "Keywords", "Draft", "Remark"},
			Filterables: []string{"Site", "CategoryID", "Tags", "Published", "UpdatedAt"},
			Orderables:  []string{"UpdatedAt", "PublishedAt"},
			Searchables: []string{"ID", "Tags", "Title", "Alt", "Description", "Keywords", "Body"},
			Requireds:   []string{"ID", "Site", "CategoryID", "ContentType", "Body"},
			Icon:        weapi.ReadIcon("/icon/newspaper.svg"),
			Styles: []string{
				"/static/admin/css/easymde.min.css",
				"/static/admin/css/jodit.min.css",
			},
			Scripts: []object.AdminScript{
				{Src: "/static/admin/js/cms_widget.js"},
				{Src: "/static/admin/js/easymde.min.js"},
				{Src: "/static/admin/js/jodit.min.js"},
				{Src: "/static/admin/js/cms_page.js", Onload: true}},
			Attributes: map[string]object.AdminAttribute{
				"ContentType": {Choices: weapi.EnabledPageContentTypes, Default: enum.ContentTypeHtml},
				"Draft":       {Default: "Your content ..."},
				"IsDraft":     {Widget: "is-draft"},
				"Published":   {Widget: "is-published"},
				"Tags":        {Widget: "tags", FilterWidget: "tags"},
				"CategoryID":  {Widget: "category-id-and-path", FilterWidget: "category-id-and-path"},
				"ID":          {Help: "ID must be unique,recommend use title slug eg: hello-world-2023"},
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
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						return page.HandleSaveDraft(c, obj)
					},
				},
				{
					Path: "duplicate",
					Name: "Duplicate",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						return page.HandleMakePageDuplicate(c, obj)
					},
				},
				{
					Path: "make_publish",
					Name: "Make Publish",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						return page.HandleMakePagePublish(c, obj, true)
					},
				},
				{
					Path: "make_un_publish",
					Name: "Make UnPublish",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						return page.HandleMakePagePublish(c, obj, false)
					},
				},
				{
					WithoutObject: true,
					Path:          "tags",
					Name:          "Query All Tags",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						return page.HandleQueryTags(c, obj, "posts")
					},
				},
			},
			BeforeCreate: func(ctx *fiber.Ctx, vptr any) error {
				post := vptr.(*model.Post)
				if post.ContentType == "" {
					post.ContentType = enum.ContentTypeMarkdown
				}
				post.Creator = *user.CurrentUser(ctx)
				post.IsDraft = true
				return nil
			},
			BeforeUpdate: func(ctx *fiber.Ctx, vptr any, vals map[string]any) error {
				post := vptr.(*model.Post)
				post.IsDraft = true
				if _, ok := vals["published"]; ok {
					post.Published = vals["published"].(bool)
					if post.Published {
						post.Body = post.Draft
						post.IsDraft = false
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
