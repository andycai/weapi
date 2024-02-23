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

	r.Get("/", handleDashboard)

	r.Post("/json", func(c *fiber.Ctx) error {
		return handleJson(c, adminObjects)
	})

	r.Post("/summary", HandleAdminSummary)
	r.Post("/tags/:content_type", handleGetTags)
}

func initAdminObject() []model.AdminObject {
	return []model.AdminObject{
		{
			Model: &model.Site{},
			Group: "Contents",
			Name:  "Site",
			Shows: []string{"Domain", "Name", "Preview", "Disallow", "UpdatedAt", "CreatedAt"},
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
	}
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAdminObject(initAdminObject())
}
