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

func initPublicNoCheckRouter(r fiber.Router) {
	auth := r.Group("/auth")
	{
		auth.Get("/login", handleSigin)
		auth.Post("/login", handleSiginAction)
		auth.Get("/register", handleSigup)
		auth.Post("/register", handleSigupAction)
		auth.Get("/logout", handleLogoutAction)
	}
}

func initAdminCheckRouter(r fiber.Router) {
}

func initAdminObject() []model.AdminObject {
	return []model.AdminObject{
		{
			Model:       &model.User{},
			Group:       "Settings",
			Name:        "User",
			Desc:        "Builtin user management system",
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
			Weight: 91,
		},
		{
			Model:       &model.Group{},
			Group:       "Settings",
			Name:        "Group",
			Desc:        "A group describes a group of users. One user can be part of many groups and one group can have many users", //
			Shows:       []string{"ID", "Name", "Extra", "UpdatedAt", "CreatedAt"},
			Editables:   []string{"Name", "Extra", "UpdatedAt"},
			Orderables:  []string{"UpdatedAt"},
			Searchables: []string{"Name"},
			Requireds:   []string{"Name"},
			Icon:        weapi.ReadIcon("/icon/group.svg"),
			AccessCheck: SuperAccessCheck,
			Weight:      92,
		},
		{
			Model:       &model.GroupMember{},
			Group:       "Settings",
			Name:        "GroupMember",
			Desc:        "Group members", //
			Shows:       []string{"ID", "User", "Group", "Role", "CreatedAt"},
			Filterables: []string{"Group", "Role", "CreatedAt"},
			Editables:   []string{"User", "Group", "Role"},
			Orderables:  []string{"CreatedAt"},
			Searchables: []string{"User", "Group"},
			Requireds:   []string{"User", "Group", "Role"},
			Icon:        weapi.ReadIcon("/icon/members.svg"),
			AccessCheck: SuperAccessCheck,
			Attributes: map[string]model.AdminAttribute{
				"Role": {
					Default: model.GroupRoleMember,
					Choices: []model.AdminSelectOption{{Label: "Admin", Value: model.GroupRoleAdmin}, {Label: "Member", Value: model.GroupRoleMember}},
				},
			},
			Weight: 93,
		},
		{
			Model:       &model.Config{},
			Group:       "Settings",
			Name:        "Config",
			Desc:        "System config with database backend, You can change it in admin page, and it will take effect immediately without restarting the server", //
			Shows:       []string{"Key", "Value", "Desc"},
			Editables:   []string{"Key", "Value", "Desc"},
			Orderables:  []string{"Key"},
			Searchables: []string{"Key", "Value", "Desc"},
			Requireds:   []string{"Key", "Value"},
			Icon:        weapi.ReadIcon("/icon/config.svg"),
			AccessCheck: SuperAccessCheck,
			Weight:      94,
		},
	}
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterPublicNoCheckRouter(KeyNoCheckRouter, initPublicNoCheckRouter)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAdminObject(initAdminObject())
	core.RegisterAdminAuth(WithAdminAuth)
}
