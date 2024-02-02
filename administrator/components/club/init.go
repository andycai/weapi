package club

import (
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.club.gorm.db"
	KeyNoCheckRouter = "admin.club.router.nocheck"
	KeyCheckRouter   = "admin.club.router.check"
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
			Model:       &model.Club{},
			Group:       "Activities",
			Name:        "Club",
			Desc:        "The club data of the website can only be in JSON format",
			Shows:       []string{"ID", "Name", "Level", "CreatedAt"},
			Editables:   []string{"Name", "Level", "Logo", "Notice", "Addr"},
			Filterables: []string{"UpdatedAt"},
			Orderables:  []string{"UpdatedAt"},
			Searchables: []string{"name"},
			Requireds:   []string{"Name"},
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
			Weight: 22,
		},
	}
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAdminObject(initAdminObject())
}
