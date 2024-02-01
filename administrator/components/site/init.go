package site

import (
	"github.com/andycai/weapi"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
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

func initAdminCheckRouter(r fiber.Router) {
	adminObjects := BuildAdminObjects(r, core.GetAdminObjects())

	r.Post("/json", func(c *fiber.Ctx) error {
		return JsonAction(c, adminObjects)
	})

	r.Post("/summary", HandleAdminSummary)
	r.Post("/tags/:content_type", handleGetTags)
}

func initAdminObject() []model.AdminObject {
	return []model.AdminObject{
		{
			Model:      &model.Site{},
			Group:      "Contents",
			Name:       "Site",
			PluralName: "Sites",
			Shows:      []string{"Domain", "Name", "Preview", "Disallow", "UpdatedAt", "CreatedAt"},
			Orders: []model.Order{
				{
					Name: "UpdatedAt",
					Op:   model.OrderOpDesc,
				},
			},
			Editables:   []string{"Domain", "Name", "Preview", "Disallow"},
			Filterables: []string{"Disallow"},
			Orderables:  []string{},
			Searchables: []string{"Domain", "Name", "Preview"},
			Requireds:   []string{"Domain"},
			Icon:        weapi.ReadIcon("/icon/desktop.svg"),
			Scripts: []model.AdminScript{
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
			Attributes:  map[string]model.AdminAttribute{"Items": {Widget: "category-item"}},
			Scripts: []model.AdminScript{
				{Src: "/static/admin/js/cms_widget.js"},
				{Src: "/static/admin/js/cms_category.js", Onload: true},
			},
			Actions: []model.AdminAction{
				{
					WithoutObject: true,
					Path:          "query_with_count",
					Name:          "Query with item count",
					Handler:       HandleQueryCategoryWithCount,
				},
			},
			Weight: 11,
		},
	}
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAdminObject(initAdminObject())
}
