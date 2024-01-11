package main

import (
	"path/filepath"

	_ "github.com/andycai/weapi/administrator/components"
	_ "github.com/andycai/weapi/components"
	"github.com/andycai/weapi/conf"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/library/authentication"
	"github.com/andycai/weapi/library/database"
	"github.com/andycai/weapi/library/renderer"
	"github.com/andycai/weapi/log"
	"github.com/andycai/weapi/middlewares"
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

	addr := viper.GetString("httpserver.addr")
	debug := viper.GetBool("debug")
	dbtype := viper.GetString("db.type")
	dsn := viper.GetString("db.dsn")

	RunSetup(addr)

	// database open and init
	db, err := database.InitRDBMS(dbtype, dsn,
		viper.GetInt("db.active"),
		viper.GetInt("db.idle"),
		viper.GetInt("db.idletimeout"))

	if err != nil {
		panic(err)
	}

	if debug {
		db = db.Debug()
	}

	dbs := []*gorm.DB{db}
	err = core.AutoMultiMigrate(dbs)
	if err != nil {
		panic(err)
	}
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
