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
			Shows:       []string{"ID", "Site", "Creator", "Club", "Name", "BeginAt", "EndAt"},
			Editables:   []string{"Site", "Club", "Creator", "Name", "Description", "BeginAt", "EndAt", "Kind", "Type", "Quota", "Waiting", "Stage", "FeeType", "Ahead", "Address"},
			Filterables: []string{"UpdatedAt", "BeginAt", "EndAt"},
			Orderables:  []string{"UpdatedAt"},
			Searchables: []string{"Name", "Description"},
			Requireds:   []string{"Name", "Description", "BeginAt", "EndAt"},
			// Icon:        weapi.ReadIcon("/icon/piece.svg"),
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
