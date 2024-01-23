package media

import (
	"github.com/andycai/weapi"
	"github.com/andycai/weapi/administrator/components/user"
	"github.com/andycai/weapi/conf"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/model"
	"github.com/andycai/weapi/object"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.media.gorm.db"
	KeyNoCheckRouter = "admin.media.router.nocheck"
	KeyCheckRouter   = "admin.media.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initCheckRouter(r fiber.Router) {
	// r.Get("/posts/manager", ManagerPage)
}

func initAdminObject() []object.AdminObject {
	return []object.AdminObject{
		{
			Model:       &model.Media{},
			Group:       "Contents",
			Name:        "Media",
			PluralName:  "Media",
			Desc:        "All kinds of media files, such as images, videos, audios, etc.",
			Shows:       []string{"Name", "ContentType", "Author", "Published", "Size", "Dimensions", "UpdatedAt"},
			Editables:   []string{"External", "PublicUrl", "Author", "Published", "PublishedAt", "Title", "Alt", "Description", "Keywords", "ContentType", "Size", "Path", "Name", "Dimensions", "StorePath", "UpdatedAt", "Ext", "Size", "StorePath", "Remark"},
			Filterables: []string{"Published", "UpdatedAt", "ContentType", "External"},
			Orderables:  []string{"UpdatedAt", "PublishedAt", "Size"},
			Searchables: []string{"Title", "Alt", "Description", "Keywords", "Path", "Path", "Name", "StorePath"},
			Requireds:   []string{"ContentType", "Size", "Path", "Name", "Dimensions", "StorePath"},
			Icon:        weapi.ReadIcon("/icon/image.svg"),
			Attributes: map[string]object.AdminAttribute{
				"ContentType": {Choices: weapi.ContentTypes},
				"Size":        {Widget: "humanize-size"},
				"Site":        {SingleChoice: true},
			},
			Scripts: []object.AdminScript{
				{Src: "/static/admin/js/cms_widget.js"},
				{Src: "/static/admin/js/cms_media.js", Onload: true},
			},
			EditPage: "./edit_media.html",
			Orders: []object.Order{
				{
					Name: "UpdatedAt",
					Op:   object.OrderOpAsc,
				},
			},
			BeforeRender: func(ctx *fiber.Ctx, vptr any) (any, error) {
				media := vptr.(*model.Media)
				mediaHost := conf.GetValue(db, enum.KEY_CMS_MEDIA_HOST)
				mediaPrefix := conf.GetValue(db, enum.KEY_CMS_MEDIA_PREFIX)
				media.BuildPublicUrls(mediaHost, mediaPrefix)
				return vptr, nil
			},
			BeforeCreate: func(ctx *fiber.Ctx, vptr any) error {
				media := vptr.(*model.Media)
				media.Creator = *user.CurrentUser(ctx)
				return nil
			},
			BeforeDelete: func(ctx *fiber.Ctx, vptr any) error {
				media := vptr.(*model.Media)
				if err := RemoveFile(media.Path, media.Name); err != nil {
					// object.Warning("Delete file failed: ", media.StorePath, err)
				}
				return nil
			},
			Actions: []object.AdminAction{
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
		},
	}
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initCheckRouter)
	core.RegisterAdminObject(initAdminObject())
}
