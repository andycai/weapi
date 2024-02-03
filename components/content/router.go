package content

import (
	"github.com/andycai/weapi/components/site"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	keyDB            = "content.gorm.db"
	keyNoCheckRouter = "content.router.nocheck"
	keyCheckRouter   = "content.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initAPICheckRouter(r fiber.Router) {
	objs := []model.WebObject{
		{
			Model:        &model.Category{},
			AllowMethods: model.GET | model.QUERY,
			Name:         "category",
			Editables:    []string{"UUID", "SiteID", "Name", "Items"},
			Filterables:  []string{},
			Orderables:   []string{},
			Searchables:  []string{"UUID", "Name", "Items"},
		},
		{
			Model:        &model.Page{},
			AllowMethods: model.GET | model.QUERY,
			Name:         "page",
			Filterables:  []string{"SiteID", "CategoryID", "CategoryPath", "Tags", "IsDraft", "Published", "ContentType"},
			Searchables:  []string{"Title", "Description", "Body"},
			Orderables:   []string{"CreatedAt", "UpdatedAt"},
			BeforeRender: BeforeRenderPage,
		},
		{
			Model:             &model.Post{},
			AllowMethods:      model.GET | model.QUERY,
			Name:              "post",
			Filterables:       []string{"SiteID", "CategoryID", "CategoryPath", "Tags", "IsDraft", "Published", "ContentType"},
			Searchables:       []string{"Title", "Description", "Body"},
			Orderables:        []string{"CreatedAt", "UpdatedAt"},
			BeforeRender:      BeforeRenderPost,
			BeforeQueryRender: BeforeQueryRenderPost,
		},
	}
	site.RegisterObjects(r, objs)
}

func init() {
	core.RegisterDatabase(keyDB, initDB)
	core.RegisterAPICheckRouter(keyNoCheckRouter, initAPICheckRouter)
}
