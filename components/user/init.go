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
}

func initCheckRouter(r fiber.Router) {
	//
}

func init() {
	core.RegisterDatabase(KeyUserDB, initDB)
	core.RegisterAPINoCheckRouter(KeyUserCheckRouter, initNoCheckRouter)
	core.RegisterAPICheckRouter(KeyUserCheckRouter, initCheckRouter)
}
