package user

import (
	"github.com/andycai/weapi/core"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.user.gorm.db"
	KeyNoCheckRouter = "admin.user.router.nocheck"
	KeyCheckRouter   = "admin.user.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initNoCheckRouter(r fiber.Router) {
	admin := r.Group("/auth")
	{
		admin.Get("/login", SigninPage)
		admin.Post("/login", SigninAction)
		admin.Get("/register", SigninPage)
		admin.Post("/register", SigninAction)
	}
}

func initCheckRouter(r fiber.Router) {
	r.Get("/logout", LogoutAction)
	r.Get("/dashboard", DashBoardPage)

	r.Get("/users/profile", ProfilePage)
	r.Post("/users/profile", ProfileSave)
	r.Get("/users/security", SecurityPage)
	r.Post("/users/password", PasswordSave)
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterNoCheckRouter(KeyNoCheckRouter, initNoCheckRouter)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initCheckRouter)
}
