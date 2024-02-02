package main

import (
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/andycai/weapi"
	_ "github.com/andycai/weapi/administrator/components"
	"github.com/andycai/weapi/administrator/components/config"
	"github.com/andycai/weapi/administrator/components/user"
	_ "github.com/andycai/weapi/components"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/enum"
	"github.com/andycai/weapi/lib/authentication"
	"github.com/andycai/weapi/lib/database"
	"github.com/andycai/weapi/log"
	"github.com/andycai/weapi/middlewares"
	"github.com/andycai/weapi/model"
	"github.com/andycai/weapi/utils/date"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

func main() {
	engine := core.ViewEngineStart()
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	var addr string
	var logFile string = config.GetEnv(enum.ENV_LOG_FILE)
	var dbDriver string = config.GetEnv(enum.ENV_DB_DRIVER)
	var dsn string = config.GetEnv(enum.ENV_DSN)
	var debug bool = config.GetEnv(enum.ENV_DEBUG) != ""
	var staticDir string = config.GetEnv(enum.ENV_STATIC_DIR)
	var htmlDir string = config.GetEnv(enum.ENV_HTML_DIR)
	var dbActive = config.GetIntEnv(enum.ENV_DB_ACTIVE)
	var dbIdle = config.GetIntEnv(enum.ENV_DB_IDLE)
	var dbTimeout = config.GetIntEnv(enum.ENV_DB_TIMEOUT)
	var lang = config.GetEnv(enum.ENV_LANG)
	var zoneOffset = config.GetIntEnv(enum.ENV_ZONE_OFFSET)

	var superUserEmail string
	var superUserPassword string

	log.Setup(debug, logFile)

	flag.StringVar(&superUserEmail, "superuser", "", "Create an super user with email")
	flag.StringVar(&superUserPassword, "password", "", "Super user password")
	flag.StringVar(&addr, "addr", ":8888", "http service address")
	flag.StringVar(&logFile, "log", logFile, "log file path")
	flag.StringVar(&dbDriver, "db", dbDriver, "database driver")
	flag.StringVar(&dsn, "dsn", dsn, "database dsn")
	flag.BoolVar(&debug, "debug", debug, "debug mode")
	flag.StringVar(&staticDir, "static", staticDir, "static file directory")
	flag.StringVar(&htmlDir, "html", htmlDir, "html file directory")
	flag.StringVar(&lang, "en", lang, "language")
	flag.IntVar(&zoneOffset, "zone", zoneOffset, "timezone offset")
	flag.IntVar(&dbActive, "db-active", dbActive, "database max active connection")
	flag.IntVar(&dbIdle, "db-idle", dbIdle, "database max idle connection")
	flag.IntVar(&dbTimeout, "db-timeout", dbTimeout, "database connection timeout")

	// setup
	if dsn == "" {
		RunSetup(addr)
	}
	var httplw io.Writer = os.Stdout
	var lw io.Writer = os.Stdout

	if !debug {
		lw, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Infof("error opening file: %v", err)
		}
		defer lw.Close()

		httplw, err := os.OpenFile("http-"+logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Infof("error opening file: %v", err)
		}
		defer httplw.Close()
	}

	// database open and init
	db, err := database.InitRDBMS(lw, dbDriver, dsn, dbActive, dbIdle, dbTimeout)

	if err != nil {
		panic(err)
	}

	if debug {
		db = db.Debug()
		db.DB()
	}

	dbs := []*gorm.DB{db}
	err = weapi.AutoMultiMigrate(dbs)
	if err != nil {
		panic(err)
	}

	core.SetupDatabase(dbs)
	sqlDb, _ := db.DB()
	authentication.SessionSetup(dbDriver, sqlDb, dsn, "auth_session")

	if superUserEmail != "" && superUserPassword != "" {
		// create super user
		err, userVo := user.GetByEmail(superUserEmail)
		if err == nil && userVo != nil {
			user.UpdatePassword(userVo, superUserPassword)
			log.Infof("password of the super user %s has been updated", superUserEmail)
		} else {
			userVo = &model.User{}
			userVo.Email = superUserEmail
			userVo.Password = core.HashPassword(superUserPassword)
			userVo.IsStaff = true
			userVo.Activated = true
			userVo.Enabled = true
			userVo.IsSuperUser = true
			err := user.CreateUser(userVo)
			if err != nil {
				panic(err)
			}
		}
	}

	date.SetZoneOffset(zoneOffset)
	core.SetLang(lang)

	// Middleware
	middlewares.Use(app, httplw)

	app.Static("/static", filepath.Join("", staticDir))
	app.Static("/admin", filepath.Join("", htmlDir))

	// router
	core.SetupRouter(app)

	// check config
	config.CheckConfig()

	log.Infof("Server started on %s", addr)
	err = app.Listen(addr)
	if err != nil {
		panic(err)
	}
	defer func() {
		//
	}()
}
