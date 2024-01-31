package user

import (
	"github.com/andycai/weapi"
	"github.com/andycai/weapi/components/user"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/andycai/weapi/object"
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
		auth.Get("/login", SigninPage)
		auth.Post("/login", SigninAction)
		auth.Get("/register", SigninPage)
		auth.Post("/register", SigninAction)
		auth.Get("/logout", LogoutAction)
	}
}

func initCheckRouter(r fiber.Router) {
}

func initAdminObject() []object.AdminObject {
	return []object.AdminObject{
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
			Orders:      []object.Order{{Name: "UpdatedAt", Op: object.OrderOpDesc}},
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
				if err, dbUser := user.Dao.GetByEmail(userVo.Email); err == nil {
					if dbUser.Password != userVo.Password {
						userVo.Password = core.HashPassword(userVo.Password)
					}
				}
				return nil
			},
			Actions: []object.AdminAction{
				{
					Path:  "toggle_enabled",
					Name:  "Toggle enabled",
					Label: "Toggle user enabled/disabled",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						userVo := obj.(*model.User)
						err := user.Dao.UpdateFields(userVo, map[string]any{"Enabled": !userVo.Enabled})
						return userVo.Enabled, err
					},
				},
				{
					Path:  "toggle_staff",
					Name:  "Toggle staff",
					Label: "Toggle user is staff or not",
					Handler: func(c *fiber.Ctx, obj any) (any, error) {
						userVo := obj.(*model.User)
						err := user.Dao.UpdateFields(userVo, map[string]any{"IsStaff": !userVo.IsStaff})
						return userVo.IsStaff, err
					},
				},
			},
			Attributes: map[string]object.AdminAttribute{
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
