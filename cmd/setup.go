package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/andycai/weapi"
	"github.com/andycai/weapi/administrator/components/user"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/lib/database"
	"github.com/andycai/weapi/middlewares"
	"github.com/andycai/weapi/utils/date"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var GitCommit string = ""
var BuildTime string = ""
var setupDoneFlag string = ".weapi_setup_done"

type SetupDBForm struct {
	Driver   string `json:"dbDriver"`
	Host     string `json:"dbHost"`
	Port     string `json:"dbPort"`
	Name     string `json:"dbName"`
	Filename string `json:"dbFilename"`
	Charset  string `json:"dbCharset"`
	User     string `json:"dbUser"`
	Password string `json:"dbPassword"`
}

type SetupSuperUserForm struct {
	DBConfig SetupDBForm `json:"dbConfig" binding:"required"`
	Username string      `json:"superUsername" binding:"required"`
	Password string      `json:"superPassword" binding:"required"`
}

type SetupSuperEnvForm struct {
	DBConfig     SetupDBForm `json:"dbConfig" binding:"required"`
	Salt         string      `json:"salt" binding:"required"`
	CookieSecret string      `json:"cookieSecret" binding:"required"`
	LogFile      string      `json:"logFile"`
}

func (f *SetupDBForm) DSN() string {
	if f.Driver == "sqlite" {
		return fmt.Sprintf("file:%s", f.Filename)
	}

	pwd := f.Password
	if pwd == "" {
		pwd = "''"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		f.User, pwd, f.Host, f.Port, f.Name, f.Charset)
}

func ok(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": true,
	})
}

func fail(c *fiber.Ctx, data string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"data": data,
	})
}

func RunSetup(addr string) {
	if _, err := os.Stat(setupDoneFlag); err != nil {
		runSetupMode(addr)
	}
}

func runSetupMode(addr string) {
	// var err error
	// log.Infof("Run setup mode")
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// log.Infof("Please visit http://%s/setup to complete install", addr)

	engine := core.ViewEngineStart()
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Middleware
	middlewares.Use(app)
	app.Static("/static", filepath.Join("", viper.GetString("app.static")))

	srv := &http.Server{Handler: adaptor.FiberApp(app)}

	app.Get("/admin", func(c *fiber.Ctx) error {
		return c.Redirect("/setup/", http.StatusFound)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/setup/")
	})

	app.Get("/setup", func(c *fiber.Ctx) error {
		osVersion := fmt.Sprintf("%s-%s (%s)", runtime.GOOS, runtime.GOARCH, runtime.Version())
		// current working directory
		cwd, _ := os.Getwd()
		data := map[string]any{
			"buildTime":    date.Format(date.Now(), "2006-01-02 15:04:05"),
			"gitCommit":    GitCommit,
			"osVersion":    osVersion,
			"cwd":          cwd,
			"enableSqlite": enableSqlite,
		}
		return c.Render("setup", data)
	})

	app.Post("/setup/ping_database", func(c *fiber.Ctx) error {
		var form SetupDBForm
		if err := c.BodyParser(&form); err != nil {
			return fail(c, err.Error())
		}

		var err error
		var db *gorm.DB
		defer func() {
			if db != nil {
				h, _ := db.DB()
				h.Close()
			}
		}()

		// log.Infof("DSN: %s", form.DSN())
		db, err = database.InitRDBMS(form.Driver, form.DSN(), 32, 30, 14400)
		if err != nil {
			return fail(c, err.Error())
		}

		return ok(c)
	})

	app.Post("/setup/migrate_database", func(c *fiber.Ctx) error {
		var form SetupDBForm
		if err := c.BodyParser(&form); err != nil {
			return fail(c, err.Error())
		}

		var err error
		var db *gorm.DB
		defer func() {
			if db != nil {
				h, _ := db.DB()
				h.Close()
			}
		}()
		db, err = database.InitRDBMS(form.Driver, form.DSN(), 32, 30, 14400)
		if err != nil {
			return fail(c, err.Error())
		}

		err = weapi.AutoMigrate(db)
		if err != nil {
			return fail(c, err.Error())
		}

		return ok(c)
	})

	app.Post("/setup/write_env", func(c *fiber.Ctx) error {
		var form SetupSuperEnvForm
		if err := c.BodyParser(&form); err != nil {
			return fail(c, err.Error())
		}
		envFile := ".env"

		lines := []string{
			fmt.Sprintf("%s=%s", "PASSWORD_SALT", form.Salt),
			fmt.Sprintf("%s=%s", "SESSION_SECRET", form.CookieSecret),
			fmt.Sprintf("LOG_FILE=%s", form.LogFile),
			fmt.Sprintf("DSN=%s", form.DBConfig.DSN()),
			fmt.Sprintf("DB_DRIVER=%s", form.DBConfig.Driver),
			fmt.Sprintf("DB_ACTIVE=%d", 32),
			fmt.Sprintf("DB_IDLE=%d", 30),
			fmt.Sprintf("DB_TIMEOUT=%d", 14400),
			fmt.Sprintf("DEBUG=%s", "true"),
			fmt.Sprintf("LANG=%s", "en"),
			fmt.Sprintf("ZONE_OFFSET=%d", 8),
			fmt.Sprintf("STATIC_DIR=%s", "static"),
			fmt.Sprintf("HTML_DIR=%s", "templates/admin"),
			fmt.Sprintf("LOG_DIR=%s", "log"),
			fmt.Sprintf("CACHE_DIR=%s", "cache"),
			fmt.Sprintf("REDIS_ADDR=%s", "127.0.0.1:6379"),
			fmt.Sprintf("REDIS_PASSWORD=%s", "i18n!@"),
			fmt.Sprintf("REDIS_DB=%d", 0),
		}
		data := strings.Join(lines, "\n") + "\n"
		if _, err := os.Stat(envFile); err == nil {
			fileData, _ := os.ReadFile(envFile)
			if fileData != nil {
				data = string(fileData) + data
			}
		}
		os.WriteFile(envFile, []byte(data), 0644)

		return ok(c)
	})

	app.Post("/setup/create_superuser", func(c *fiber.Ctx) error {
		var form SetupSuperUserForm
		if err := c.BodyParser(&form); err != nil {
			return fail(c, err.Error())
		}

		var err error
		var db *gorm.DB
		defer func() {
			if db != nil {
				h, _ := db.DB()
				h.Close()
			}
		}()

		db, err = database.InitRDBMS(form.DBConfig.Driver, form.DBConfig.DSN(), 32, 30, 14400)
		if err != nil {
			return fail(c, err.Error())
		}
		err = weapi.AutoMigrate(db)
		if err != nil {
			return fail(c, err.Error())
		}

		core.SetupDatabase([]*gorm.DB{db})

		err, u := user.GetByEmail(form.Username)
		if err == nil && u != nil {
			user.UpdatePassword(u, form.Password)
		} else {
			err = user.CreateUser(u)
			if err != nil {
				panic(err)
			}
		}
		u.IsStaff = true
		u.Activated = true
		u.Enabled = true
		u.IsSuperUser = true
		db.Save(u)
		// log.Infof("Create super user: %s", form.Username)

		return ok(c)
	})

	app.Post("/setup/restart", func(c *fiber.Ctx) error {
		os.WriteFile(setupDoneFlag, []byte("done"), 0644)
		time.AfterFunc(500*time.Millisecond, func() {
			// log.Infof("Restarting...")
			srv.Shutdown(context.Background())
		})

		return ok(c)
	})
	// err = app.Listen(viper.GetString("httpserver.addr"))
	// if err != nil {
	// 	panic(err)
	// }
	srv.Serve(ln)
}
