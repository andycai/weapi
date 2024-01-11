package post

import (
	"github.com/andycai/weapi/core"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.post.gorm.db"
	KeyNoCheckRouter = "admin.post.router.nocheck"
	KeyCheckRouter   = "admin.post.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initCheckRouter(r fiber.Router) {
	r.Get("/posts/manager", ManagerPage)
	r.Get("/posts/editor", EditorPage)
	r.Get("/posts/editor/:id", EditorPage)
	r.Post("/posts/editor", Create)
	r.Post("/posts/editor/:id", Update)
	r.Get("/posts/movetotrash/:id", MoveToTrashByID)
	r.Post("/posts/movetotrash", MoveToTrash)
	r.Get("/posts/restore/:id", RestoreByID)
	r.Post("/posts/restore", Restore)

	r.Get("/categories/manager", ManagerCategoryPage)
	r.Get("/categories/editor", EditorCategoryPage)
	r.Get("/categories/editor/:id", EditorCategoryPage)
	r.Post("/categories/editor", CreateCategory)
	r.Post("/categories/editor/:id", UpdateCategory)
	r.Post("/categories/delete", DeleteCategories)

	r.Get("/tags/manager", ManagerTagsPage)
	r.Get("/tags/editor", EditorTagPage)
	r.Post("/tags/editor", CreateTag)
	r.Post("/tags/delete", DeleteTags)
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initCheckRouter)
}
