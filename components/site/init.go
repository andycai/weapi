package site

import (
	"github.com/andycai/weapi/components/page"
	"github.com/andycai/weapi/components/post"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyPageDB            = "entity.gorm.db"
	KeyPageNoCheckRouter = "entity.router.nocheck"
	KeyPageCheckRouter   = "entity.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initNoCheckRouter(r fiber.Router) {

}

func initCheckRouter(r fiber.Router) {
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
			GetDB:        post.GetPostOrPageDB,
			BeforeRender: page.BeforeRenderPage,
		},
		{
			Model:             &model.Post{},
			AllowMethods:      model.GET | model.QUERY,
			Name:              "post",
			Filterables:       []string{"SiteID", "CategoryID", "CategoryPath", "Tags", "IsDraft", "Published", "ContentType"},
			Searchables:       []string{"Title", "Description", "Body"},
			Orderables:        []string{"CreatedAt", "UpdatedAt"},
			GetDB:             post.GetPostOrPageDB,
			BeforeRender:      post.BeforeRenderPost,
			BeforeQueryRender: post.BeforeQueryRenderPost,
		},
	}
	RegisterObjects(r, objs)
}

func init() {
	core.RegisterDatabase(KeyPageDB, initDB)
	core.RegisterAPINoCheckRouter(KeyPageNoCheckRouter, initNoCheckRouter)
	core.RegisterAPICheckRouter(KeyPageNoCheckRouter, initCheckRouter)
}
