package content

import (
	"github.com/andycai/weapi"
	"github.com/andycai/weapi/administrator/components/user"
	userapi "github.com/andycai/weapi/components/user"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/log"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.content.gorm.db"
	KeyNoCheckRouter = "admin.content.router.nocheck"
	KeyCheckRouter   = "admin.content.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initPublicNoCheckRouter(r fiber.Router) {
	mediaPrefix := user.GetValue(enum.KEY_CMS_MEDIA_PREFIX)
	if mediaPrefix == "" {
		mediaPrefix = "/media/"
	}
	g := r.Group(mediaPrefix, userapi.WithAPIAuth)
	g.Get("/*", handleMedia)
}

func initAdminCheckRouter(r fiber.Router) {
}

func initAdminObject() []model.AdminObject {
	return []model.AdminObject{
		{
			Model:       &model.Category{},
			Group:       "Contents",
			Name:        "Category",
			Desc:        "The category of articles and pages can be multi-level",
			Shows:       []string{"Name", "UUID", "Site", "Items"},
			Editables:   []string{"Name", "UUID", "Site", "Items"},
			Orderables:  []string{},
			Searchables: []string{"UUID", "Site", "Items", "Name"},
			Requireds:   []string{"UUID", "Site", "Items", "Name"},
			Icon:        weapi.ReadIcon("/icon/swatch.svg"),
			Attributes:  map[string]model.AdminAttribute{"Items": {Widget: "category-item"}},
			Scripts: []model.AdminScript{
				{Src: "/static/admin/js/cms_widget.js"},
				{Src: "/static/admin/js/cms_category.js", Onload: true},
			},
			Actions: []model.AdminAction{
				{
					WithoutObject: true,
					Path:          "query_with_count",
					Name:          "Query with item count",
					Handler:       handleQueryCategoryWithCount,
				},
			},
			Weight: 11,
		},
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
			Icon:        weapi.ReadIcon("/icon/piece.svg"),
			Styles: []string{
				"/static/admin/css/jsoneditor-9.10.2.min.css",
			},
			Scripts: []model.AdminScript{
				{Src: "/static/admin/js/cms_widget.js"},
				{Src: "/static/admin/js/jsoneditor-9.10.2.min.js"},
				{Src: "/static/admin/js/cms_page.js", Onload: true}},
			Attributes: map[string]model.AdminAttribute{
				"ContentType": {Choices: weapi.EnabledPageContentTypes, Default: enum.ContentTypeJson},
				"Draft":       {Default: "{}"},
				"IsDraft":     {Widget: "is-draft"},
				"Published":   {Widget: "is-published"},
				"Tags":        {Widget: "tags", FilterWidget: "tags"},
				"CategoryID":  {Widget: "category-id-and-path", FilterWidget: "category-id-and-path"},
				"ID":          {Help: "ID must be unique,recommend use page url eg: about-us"},
			},
			EditPage: "/html/edit_page.html",
			Orders: []model.Order{
				{
					Name: "UpdatedAt",
					Op:   model.OrderOpDesc,
				},
			},
			Actions: []model.AdminAction{
				{
					WithoutObject: true,
					Path:          "save_draft",
					Name:          "Safe Draft",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						return handleSaveDraft(c, obj)
					},
				},
				{
					Path: "duplicate",
					Name: "Duplicate",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						return handleMakePageDuplicate(c, obj)
					},
				},
				{
					Path: "make_publish",
					Name: "Make Publish",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						return handleMakePagePublish(c, obj, true)
					},
				},
				{
					Path: "make_un_publish",
					Name: "Make UnPublish",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						return handleMakePagePublish(c, obj, false)
					},
				},
				{
					WithoutObject: true,
					Path:          "tags",
					Name:          "Query All Tags",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						return handleQueryPageTags(c, obj, "pages")
					},
				},
			},
			BeforeCreate: func(ctx *fiber.Ctx, vptr any) error {
				page := vptr.(*model.Page)
				page.ContentType = enum.ContentTypeJson
				page.Creator = *user.Current(ctx)
				page.IsDraft = true
				return nil
			},
			BeforeUpdate: func(ctx *fiber.Ctx, vptr any, vals map[string]any) error {
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
			Weight: 12,
		},
		{
			Model:       &model.Post{},
			Group:       "Contents",
			Name:        "Post",
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
			Scripts: []model.AdminScript{
				{Src: "/static/admin/js/cms_widget.js"},
				{Src: "/static/admin/js/easymde.min.js"},
				{Src: "/static/admin/js/jodit.min.js"},
				{Src: "/static/admin/js/cms_page.js", Onload: true}},
			Attributes: map[string]model.AdminAttribute{
				"ContentType": {Choices: weapi.EnabledPageContentTypes, Default: enum.ContentTypeHtml},
				"Draft":       {Default: "Your content ..."},
				"IsDraft":     {Widget: "is-draft"},
				"Published":   {Widget: "is-published"},
				"Tags":        {Widget: "tags", FilterWidget: "tags"},
				"CategoryID":  {Widget: "category-id-and-path", FilterWidget: "category-id-and-path"},
				"ID":          {Help: "ID must be unique,recommend use title slug eg: hello-world-2023"},
			},
			EditPage: "/html/edit_page.html",
			Orders: []model.Order{
				{
					Name: "UpdatedAt",
					Op:   model.OrderOpDesc,
				},
			},
			Actions: []model.AdminAction{
				{
					WithoutObject: true,
					Path:          "save_draft",
					Name:          "Safe Draft",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						return handleSaveDraft(c, obj)
					},
				},
				{
					Path: "duplicate",
					Name: "Duplicate",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						return handleMakePageDuplicate(c, obj)
					},
				},
				{
					Path: "make_publish",
					Name: "Make Publish",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						return handleMakePagePublish(c, obj, true)
					},
				},
				{
					Path: "make_un_publish",
					Name: "Make UnPublish",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						return handleMakePagePublish(c, obj, false)
					},
				},
				{
					WithoutObject: true,
					Path:          "tags",
					Name:          "Query All Tags",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						return handleQueryPostTags(c, obj, "posts")
					},
				},
			},
			BeforeCreate: func(ctx *fiber.Ctx, vptr any) error {
				post := vptr.(*model.Post)
				if post.ContentType == "" {
					post.ContentType = enum.ContentTypeMarkdown
				}
				post.Creator = *user.Current(ctx)
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
			Weight: 13,
		},
		{
			Model:       &model.Media{},
			Group:       "Contents",
			Name:        "Media",
			Desc:        "All kinds of media files, such as images, videos, audios, etc.",
			Shows:       []string{"Name", "ContentType", "Author", "Published", "Size", "Dimensions", "UpdatedAt"},
			Editables:   []string{"External", "PublicUrl", "Author", "Published", "PublishedAt", "Title", "Alt", "Description", "Keywords", "ContentType", "Size", "Path", "Name", "Dimensions", "StorePath", "UpdatedAt", "Ext", "Size", "StorePath", "Remark"},
			Filterables: []string{"Published", "UpdatedAt", "ContentType", "External"},
			Orderables:  []string{"UpdatedAt", "PublishedAt", "Size"},
			Searchables: []string{"Title", "Alt", "Description", "Keywords", "Path", "Path", "Name", "StorePath"},
			Requireds:   []string{"ContentType", "Size", "Path", "Name", "Dimensions", "StorePath"},
			Icon:        weapi.ReadIcon("/icon/image.svg"),
			Attributes: map[string]model.AdminAttribute{
				"ContentType": {Choices: weapi.ContentTypes},
				"Size":        {Widget: "humanize-size"},
				"Site":        {SingleChoice: true},
			},
			Scripts: []model.AdminScript{
				{Src: "/static/admin/js/cms_widget.js"},
				{Src: "/static/admin/js/cms_media.js", Onload: true},
			},
			EditPage: "/html/edit_media.html",
			Orders: []model.Order{
				{
					Name: "UpdatedAt",
					Op:   model.OrderOpAsc,
				},
			},
			BeforeRender: func(ctx *fiber.Ctx, vptr any) (any, error) {
				media := vptr.(*model.Media)
				mediaHost := user.GetValue(enum.KEY_CMS_MEDIA_HOST)
				mediaPrefix := user.GetValue(enum.KEY_CMS_MEDIA_PREFIX)
				media.BuildPublicUrls(mediaHost, mediaPrefix)
				return vptr, nil
			},
			BeforeCreate: func(ctx *fiber.Ctx, vptr any) error {
				media := vptr.(*model.Media)
				media.Creator = *user.Current(ctx)
				return nil
			},
			BeforeDelete: func(ctx *fiber.Ctx, vptr any) error {
				media := vptr.(*model.Media)
				if err := removeFile(media.Path, media.Name); err != nil {
					log.Infof("Delete file failed: %s, %s", media.StorePath, err)
				}
				return nil
			},
			Actions: []model.AdminAction{
				{
					Path: "make_publish",
					Name: "Make Publish",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						return handleMakeMediaPublish(c, obj, true)
					},
				},
				{
					Path: "make_un_publish",
					Name: "Make UnPublish",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						return handleMakeMediaPublish(c, obj, false)
					},
				},
				{
					WithoutObject: true,
					Path:          "folders",
					Name:          "Folders",
					Handler:       handleListFolders,
				},
				{
					WithoutObject: true,
					Path:          "new_folder",
					Name:          "New Folder",
					Handler:       handleNewFolder,
				},
				{
					WithoutObject: true,
					Path:          "upload",
					Name:          "Upload",
					Handler:       handleUpload,
				},
				{
					WithoutObject: true,
					Path:          "remove_dir",
					Name:          "Remove directory",
					Handler:       handleRemoveDirectory,
				},
			},
			Weight: 14,
		},
	}
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterPublicNoCheckRouter(KeyNoCheckRouter, initPublicNoCheckRouter)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAdminObject(initAdminObject())
}
