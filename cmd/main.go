package main

import (
	"path/filepath"

	_ "github.com/andycai/werite/administrator/components"
	_ "github.com/andycai/werite/components"
	"github.com/andycai/werite/conf"
	"github.com/andycai/werite/core"
	"github.com/andycai/werite/library/authentication"
	"github.com/andycai/werite/library/database"
	"github.com/andycai/werite/library/renderer"
	"github.com/andycai/werite/log"
	"github.com/andycai/werite/middlewares"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {
	engine := renderer.ViewEngineStart()
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	log.Setup()
	conf.ReadConf()

	// database open and init
	db, err := database.InitRDBMS(viper.GetString("db.type"),
		viper.GetString("db.dsn"),
		viper.GetInt("db.active"),
		viper.GetInt("db.idle"),
		viper.GetInt("db.idletimeout"))
	if err != nil {
		panic(err)
	}
	// dao.SetDefault(db)
	dbs := []*gorm.DB{db}
	core.AutoMigrate(dbs)
	core.SetupDatabase(dbs)
	authentication.SessionStart()
	core.SetZoneOffset(viper.GetInt("app.zoneoffset"))
	core.SetLang(viper.GetString("app.lang"))

	// Middleware
	middlewares.Use(app)

	app.Static("/static", filepath.Join("", viper.GetString("app.static")))

	// router
	core.SetupRouter(app)

	err = app.Listen(viper.GetString("httpserver.addr"))
	if err != nil {
		panic(err)
	}
	defer func() {
		//
	}()
}
