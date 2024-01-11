package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/andycai/werite/core"
	"github.com/andycai/werite/library/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"gorm.io/gorm"
)

var GitCommit string
var BuildTime string
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

func runSetupMode(addr string) {
	// carrot.Warning("Run setup mode")
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// carrot.Warning("Please visit http://", addr, "/setup to complete install")

	app := fiber.New()

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
			"buildTime":    BuildTime,
			"gitCommit":    GitCommit,
			"osVersion":    osVersion,
			"cwd":          cwd,
			"enableSqlite": enableSqlite,
		}
		return core.Render(c, "setup", data)
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

		// carrot.Warning("DSN", form.DSN())
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

		err = core.AutoMigrate(db)
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
			// fmt.Sprintf("%s=%s", carrot.ENV_SALT, form.Salt),
			// fmt.Sprintf("%s=%s", carrot.ENV_SESSION_SECRET, form.CookieSecret),
			fmt.Sprintf("LOG_FILE=%s", form.LogFile),
			fmt.Sprintf("DSN=%s", form.DBConfig.DSN()),
			fmt.Sprintf("DB_DRIVER=%s", form.DBConfig.Driver),
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
		err = core.AutoMigrate(db)
		if err != nil {
			return fail(c, err.Error())
		}

		// u, err := carrot.GetUserByEmail(db, form.Username)
		// if err == nil && u != nil {
		// 	carrot.SetPassword(db, u, form.Password)
		// 	carrot.Warning("Update super with new password")
		// } else {
		// 	u, err = carrot.CreateUser(db, form.Username, form.Password)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// }
		// u.IsStaff = true
		// u.Activated = true
		// u.Enabled = true
		// u.IsSuperUser = true
		// db.Save(u)
		// carrot.Warning("Create super user:", form.Username)

		return ok(c)
	})

	app.Post("/setup/restart", func(c *fiber.Ctx) error {
		os.WriteFile(setupDoneFlag, []byte("done"), 0644)
		time.AfterFunc(500*time.Millisecond, func() {
			// carrot.Warning("Restarting...")
			srv.Shutdown(context.Background())
		})

		return ok(c)
	})

	srv.Serve(ln)
}
