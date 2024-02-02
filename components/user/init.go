package user

import (
	"github.com/andycai/weapi/core"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	keyDB            = "user.gorm.db"
	keyNoCheckRouter = "user.router.nocheck"
	keyCheckRouter   = "user.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initAPICheckRouter(r fiber.Router) {
	//
}

func init() {
	core.RegisterDatabase(keyDB, initDB)
	core.RegisterAPICheckRouter(keyCheckRouter, initAPICheckRouter)
	core.RegisterAPIAuth(WithAPIAuth)
}
