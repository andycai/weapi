package user

import (
	"github.com/andycai/weapi/core"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyUserDB            = "user.gorm.db"
	KeyUserNoCheckRouter = "user.router.nocheck"
	KeyUserCheckRouter   = "user.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initNoCheckRouter(r fiber.Router) {
	r.Get("/sign-in", SignInPage)
	r.Get("/sign-up", SignUpPage)

	// HTMX
	r.Get("/htmx/sign-in", HTMXSignInPage)
	r.Post("/htmx/sign-in", HTMXSignInAction)
	r.Post("/htmx/sign-out", HTMXSignOut)

	r.Get("/htmx/sign-up", HTMXSignUpPage)
	r.Post("/htmx/sign-up", HTMXSignUpAction)
}

func initCheckRouter(r fiber.Router) {
	//
}

func init() {
	core.RegisterDatabase(KeyUserDB, initDB)
	core.RegisterNoCheckRouter(KeyUserCheckRouter, initNoCheckRouter)
	core.RegisterCheckRouter(KeyUserCheckRouter, initCheckRouter)
}
