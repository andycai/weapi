package activity

import (
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.activity.gorm.db"
	KeyNoCheckRouter = "admin.activity.router.nocheck"
	KeyCheckRouter   = "admin.activity.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initAdminCheckRouter(r fiber.Router) {
}

func initAdminObject() []model.AdminObject {
	return []model.AdminObject{
		{
			Model:       &model.Activity{},
			Group:       "Activities",
			Name:        "Activity",
			Desc:        "The activity data of the website can only be in JSON format",
			Shows:       []string{"ID", "Site", "Club", "Title", "BeginAt", "EndAt"},
			Editables:   []string{"Site", "Club", "Creator", "Title", "Description", "BeginAt", "EndAt", "Kind", "Type", "Quota"},
			Filterables: []string{"UpdatedAt", "BeginAt", "EndAt"},
			Orderables:  []string{"UpdatedAt"},
			Searchables: []string{"Title", "Description"},
			Requireds:   []string{"Title", "Description", "BeginAt", "EndAt"},
			// Icon:        weapi.ReadIcon("/icon/piece.svg"),
			Styles: []string{
				"/static/admin/css/jsoneditor-9.10.2.min.css",
			},
			Scripts: []model.AdminScript{
				{Src: "/static/admin/js/cms_widget.js"},
				{Src: "/static/admin/js/jsoneditor-9.10.2.min.js"},
			},
			Orders: []model.Order{
				{
					Name: "UpdatedAt",
					Op:   model.OrderOpDesc,
				},
			},
			Actions: []model.AdminAction{},
			BeforeCreate: func(ctx *fiber.Ctx, vptr any) error {
				return nil
			},
			BeforeUpdate: func(ctx *fiber.Ctx, vptr any, vals map[string]any) error {
				return nil
			},
			Weight: 21,
		},
	}
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAdminObject(initAdminObject())
}
