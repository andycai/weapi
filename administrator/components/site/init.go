package site

import (
	"github.com/andycai/weapi"
	"github.com/andycai/weapi/administrator/components/entity"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/andycai/weapi/object"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.site.gorm.db"
	KeyNoCheckRouter = "admin.site.router.nocheck"
	KeyCheckRouter   = "admin.site.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initCheckRouter(r fiber.Router) {
	// r.Get("/posts/manager", ManagerPage)
}

func initAdminObject() []object.AdminObject {
	return []object.AdminObject{
		{
			Model:      &model.Site{},
			Group:      "Contents",
			Name:       "Site",
			PluralName: "Sites",
			Shows:      []string{"Domain", "Name", "Preview", "Disallow", "UpdatedAt", "CreatedAt"},
			Orders: []object.Order{
				{
					Name: "UpdatedAt",
					Op:   object.OrderOpDesc,
				},
			},
			Editables:   []string{"Domain", "Name", "Preview", "Disallow"},
			Filterables: []string{"Disallow"},
			Orderables:  []string{},
			Searchables: []string{"Domain", "Name", "Preview"},
			Requireds:   []string{"Domain"},
			Icon:        weapi.ReadIcon("/icon/desktop.svg"),
			Scripts: []object.AdminScript{
				{Src: "/static/admin/js/cms_site.js", Onload: true},
			},
			Weight: 10,
		},
		{
			Model:       &model.Category{},
			Group:       "Contents",
			Name:        "Category",
			PluralName:  "Categories",
			Desc:        "The category of articles and pages can be multi-level",
			Shows:       []string{"Name", "UUID", "Site", "Items"},
			Editables:   []string{"Name", "UUID", "Site", "Items"},
			Orderables:  []string{},
			Searchables: []string{"UUID", "Site", "Items", "Name"},
			Requireds:   []string{"UUID", "Site", "Items", "Name"},
			Icon:        weapi.ReadIcon("/icon/swatch.svg"),
			Attributes:  map[string]object.AdminAttribute{"Items": {Widget: "category-item"}},
			Scripts: []object.AdminScript{
				{Src: "/static/admin/js/cms_widget.js"},
				{Src: "/static/admin/js/cms_category.js", Onload: true},
			},
			Actions: []object.AdminAction{
				{
					WithoutObject: true,
					Path:          "query_with_count",
					Name:          "Query with item count",
					Handler:       entity.HandleQueryCategoryWithCount,
				},
			},
			Weight: 11,
		},
	}
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initCheckRouter)
	core.RegisterAdminObject(initAdminObject())
}
