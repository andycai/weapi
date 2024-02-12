package site

import (
	"github.com/andycai/weapi/api/site"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	keyDB            = "activity.gorm.db"
	keyNoCheckRouter = "activity.router.nocheck"
	keyCheckRouter   = "activity.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initAPICheckRouter(r fiber.Router) {
	objs := []model.WebObject{
		{
			Model:        &model.Activity{},
			AllowMethods: model.GET | model.QUERY,
			Name:         "activity",
			Editables:    []string{"Name", "Description", "Quota", "Waiting", "Stage", "FeeType", "FeeMale", "FeeFemale", "Address", "Ahead", "BeginAt", "EndAt"},
			Filterables:  []string{"SiteID", "BeginAt", "EndAt"},
			Orderables:   []string{"UpdatedAt"},
			Searchables:  []string{"Name", "Description"},
		},
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
