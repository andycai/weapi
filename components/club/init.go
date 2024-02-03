package site

import (
	"github.com/andycai/weapi/components/site"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	keyDB            = "club.gorm.db"
	keyNoCheckRouter = "club.router.nocheck"
	keyCheckRouter   = "club.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initAPICheckRouter(r fiber.Router) {
	objs := []model.WebObject{
		{
			Model:        &model.Club{},
			AllowMethods: model.GET | model.QUERY,
			Name:         "club",
			Editables:    []string{"Name", "Description", "Logo", "Notice", "Address"},
			Filterables:  []string{"SiteID"},
			Orderables:   []string{"UpdatedAt"},
			Searchables:  []string{"Name", "Description"},
		},
	}
	site.RegisterObjects(r, objs)
}

func init() {
	core.RegisterDatabase(keyDB, initDB)
	core.RegisterAPICheckRouter(keyNoCheckRouter, initAPICheckRouter)
}
