package user

import (
	"github.com/andycai/weapi"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
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

func initRootNoCheckRouter(r fiber.Router) {
	auth := r.Group("/auth")
	{
		auth.Get("/login", handleSigin)
		auth.Post("/login", handleSiginAction)
		auth.Get("/register", handleSigup)
		auth.Post("/register", handleSigupAction)
		auth.Get("/logout", handleLogoutAction)
	}
}

func initCheckRouter(r fiber.Router) {
}

func initAdminObject() []model.AdminObject {
	return []model.AdminObject{
		{
			Model:       &model.User{},
			Group:       "Settings",
			Name:        "User",
			Desc:        "Builtin user management system",
			PluralName:  "Users",
			Shows:       []string{"ID", "Email", "Username", "FirstName", "ListName", "IsStaff", "IsSuperUser", "Enabled", "Activated", "UpdatedAt", "LastLogin", "LastLoginIP", "Source", "Locale", "Timezone"},
			Editables:   []string{"Email", "Password", "Username", "FirstName", "ListName", "IsStaff", "IsSuperUser", "Enabled", "Activated", "Profile", "Source", "Locale", "Timezone"},
			Filterables: []string{"CreatedAt", "UpdatedAt", "Username", "IsStaff", "IsSuperUser", "Enabled", "Activated "},
			Orderables:  []string{"CreatedAt", "UpdatedAt", "Enabled", "Activated"},
			Searchables: []string{"Username", "Email", "FirstName", "ListName"},
			Orders:      []model.Order{{Name: "UpdatedAt", Op: model.OrderOpDesc}},
			Icon:        weapi.ReadIcon("/icon/user.svg"),
			AccessCheck: SuperAccessCheck,
			BeforeCreate: func(c *fiber.Ctx, obj any) error {
				user := obj.(*model.User)
				if user.Password != "" {
					user.Password = core.HashPassword(user.Password)
				}
				user.Source = "admin"
				return nil
			},
			BeforeUpdate: func(c *fiber.Ctx, obj any, vals map[string]any) error {
				userVo := obj.(*model.User)
				if err, dbUser := GetByEmail(userVo.Email); err == nil {
					if dbUser.Password != userVo.Password {
						userVo.Password = core.HashPassword(userVo.Password)
					}
				}
				return nil
			},
			Actions: []model.AdminAction{
				{
					Path:  "toggle_enabled",
					Name:  "Toggle enabled",
					Label: "Toggle user enabled/disabled",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						userVo := obj.(*model.User)
						err := UpdateFields(userVo, map[string]any{"Enabled": !userVo.Enabled})
						return userVo.Enabled, err
					},
				},
				{
					Path:  "toggle_staff",
					Name:  "Toggle staff",
					Label: "Toggle user is staff or not",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						userVo := obj.(*model.User)
						err := UpdateFields(userVo, map[string]any{"IsStaff": !userVo.IsStaff})
						return userVo.IsStaff, err
					},
				},
			},
			Attributes: map[string]model.AdminAttribute{
				"Password": {
					Widget: "password",
				},
			},
			Weight: 21,
		},
	}
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterRootNoCheckRouter(KeyNoCheckRouter, initRootNoCheckRouter)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initCheckRouter)
	core.RegisterAdminObject(initAdminObject())
}
