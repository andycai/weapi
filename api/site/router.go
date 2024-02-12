package site

import (
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	keyDB            = "site.gorm.db"
	keyNoCheckRouter = "site.router.nocheck"
	keyCheckRouter   = "site.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initAPICheckRouter(r fiber.Router) {
	objs := []model.WebObject{
		{
			Model:        &model.Site{},
			AllowMethods: model.GET | model.QUERY,
			Name:         "site",
			Editables:    []string{"Domain", "Name", "Preview", "Disallow"},
			Filterables:  []string{},
			Orderables:   []string{},
			Searchables:  []string{"Domain", "Name"},
		},
	}
	RegisterObjects(r, objs)
}

func init() {
	core.RegisterDatabase(keyDB, initDB)
	core.RegisterAPICheckRouter(keyNoCheckRouter, initAPICheckRouter)
}
