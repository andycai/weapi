package post

import (
	"github.com/andycai/weapi/components/site"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	keyDB            = "post.gorm.db"
	keyNoCheckRouter = "post.router.nocheck"
	keyCheckRouter   = "post.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initAPICheckRouter(r fiber.Router) {
	objs := []model.WebObject{
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
	core.RegisterAPICheckRouter(keyCheckRouter, initAPICheckRouter)
}
