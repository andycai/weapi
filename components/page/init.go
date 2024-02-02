package page

import (
	"github.com/andycai/weapi/components/site"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	keyDB            = "page.gorm.db"
	keyNoCheckRouter = "page.router.nocheck"
	keyCheckRouter   = "page.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initAPICheckRouter(r fiber.Router) {
	objs := []model.WebObject{
		{
			Model:        &model.Page{},
			AllowMethods: model.GET | model.QUERY,
			Name:         "page",
			Filterables:  []string{"SiteID", "CategoryID", "CategoryPath", "Tags", "IsDraft", "Published", "ContentType"},
			Searchables:  []string{"Title", "Description", "Body"},
			Orderables:   []string{"CreatedAt", "UpdatedAt"},
			BeforeRender: BeforeRenderPage,
		},
	}
	site.RegisterObjects(r, objs)
}

func init() {
	core.RegisterDatabase(keyDB, initDB)
	core.RegisterAPICheckRouter(keyNoCheckRouter, initAPICheckRouter)
}
