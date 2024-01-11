package page

import (
	"github.com/andycai/weapi/core"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.page.gorm.db"
	KeyNoCheckRouter = "admin.page.router.nocheck"
	KeyCheckRouter   = "admin.page.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initCheckRouter(r fiber.Router) {
	r.Get("/pages/manager", ManagerPage)

	r.Get("/pages/editor", EditorPage)
	r.Get("/pages/editor/:id", EditorPage)
	r.Post("/pages/editor", Create)
	r.Post("/pages/editor/:id", Update)
	r.Get("/pages/movetotrash/:id", MoveToTrashByID)
	r.Post("/pages/movetotrash", MoveToTrash)
	r.Get("/pages/restore/:id", RestoreByID)
	r.Post("/pages/restore", Restore)
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initCheckRouter)
}
