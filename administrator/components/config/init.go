package config

import (
	"github.com/andycai/weapi"
	"github.com/andycai/weapi/administrator/components/user"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.config.gorm.db"
	KeyNoCheckRouter = "admin.config.router.nocheck"
	KeyCheckRouter   = "admin.config.router.check"
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
			Model:       &model.Config{},
			Group:       "Settings",
			Name:        "Config",
			Desc:        "System config with database backend, You can change it in admin page, and it will take effect immediately without restarting the server", //
			PluralName:  "Configs",
			Shows:       []string{"Key", "Value", "Desc"},
			Editables:   []string{"Key", "Value", "Desc"},
			Orderables:  []string{"Key"},
			Searchables: []string{"Key", "Value", "Desc"},
			Requireds:   []string{"Key", "Value"},
			Icon:        weapi.ReadIcon("/icon/config.svg"),
			AccessCheck: user.SuperAccessCheck,
			Weight:      24,
		},
	}
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAdminObject(initAdminObject())
}
