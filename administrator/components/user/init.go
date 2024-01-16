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

func initNoCheckRouter(r fiber.Router) {
	admin := r.Group("/auth")
	{
		admin.Get("/login", SigninPage)
		admin.Post("/login", SigninAction)
		admin.Get("/register", SigninPage)
		admin.Post("/register", SigninAction)
	}
}

func initCheckRouter(r fiber.Router) {
	r.Get("/logout", LogoutAction)
	r.Get("/dashboard", DashBoardPage)

	r.Get("/users/profile", ProfilePage)
	r.Post("/users/profile", ProfileSave)
	r.Get("/users/security", SecurityPage)
	r.Post("/users/password", PasswordSave)
}

func initAdminObject() []object.AdminObject {
	iconUser, _ := weapi.EmbedAssets.ReadFile("static/admin/icon_user.svg")
	iconGroup, _ := weapi.EmbedAssets.ReadFile("static/admin/icon_group.svg")
	iconMembers, _ := weapi.EmbedAssets.ReadFile("static/admin/icon_members.svg")
	iconConfig, _ := weapi.EmbedAssets.ReadFile("static/admin/icon_config.svg")

	return []object.AdminObject{
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
			Orders:      []object.Order{{Name: "UpdatedAt", Op: object.OrderOpDesc}},
			Icon:        &object.AdminIcon{SVG: string(iconUser)},
			AccessCheck: superAccessCheck,
			BeforeCreate: func(db *gorm.DB, c *fiber.Ctx, obj any) error {
				user := obj.(*model.User)
				if user.Password != "" {
					user.Password = core.HashPassword(user.Password)
				}
				user.Source = "admin"
				return nil
			},
			BeforeUpdate: func(db *gorm.DB, c *fiber.Ctx, obj any, vals map[string]any) error {
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
					Handler: func(db *gorm.DB, c *fiber.Ctx, obj any) (any, error) {
						userVo := obj.(*model.User)
						err := user.Dao.UpdateFields(userVo, map[string]any{"Enabled": !userVo.Enabled})
						return userVo.Enabled, err
					},
				},
				{
					Path:  "toggle_staff",
					Name:  "Toggle staff",
					Label: "Toggle user is staff or not",
					Handler: func(db *gorm.DB, c *fiber.Ctx, obj any) (any, error) {
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
		},
		{
			Model:       &model.Group{},
			Group:       "Settings",
			Name:        "Group",
			Desc:        "A group describes a group of users. One user can be part of many groups and one group can have many users", //
			Shows:       []string{"ID", "Name", "Extra", "UpdatedAt", "CreatedAt"},
			Editables:   []string{"ID", "Name", "Extra", "UpdatedAt"},
			Orderables:  []string{"UpdatedAt"},
			Searchables: []string{"Name"},
			Requireds:   []string{"Name"},
			Icon:        &object.AdminIcon{SVG: string(iconGroup)},
			AccessCheck: superAccessCheck,
		},
		{
			Model:       &model.GroupMember{},
			Group:       "Settings",
			Name:        "GroupMember",
			Desc:        "Group members", //
			Shows:       []string{"ID", "User", "Group", "Role", "CreatedAt"},
			Filterables: []string{"Group", "Role", "CreatedAt"},
			Editables:   []string{"ID", "User", "Group", "Role"},
			Orderables:  []string{"CreatedAt"},
			Searchables: []string{"User", "Group"},
			Requireds:   []string{"User", "Group", "Role"},
			Icon:        &object.AdminIcon{SVG: string(iconMembers)},
			AccessCheck: superAccessCheck,
			Attributes: map[string]object.AdminAttribute{
				"Role": {
					Default: model.GroupRoleMember,
					Choices: []object.AdminSelectOption{{Label: "Admin", Value: model.GroupRoleAdmin}, {Label: "Member", Value: model.GroupRoleMember}},
				},
			},
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
			Icon:        &object.AdminIcon{SVG: string(iconConfig)},
			AccessCheck: superAccessCheck,
		},
	}
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterNoCheckRouter(KeyNoCheckRouter, initNoCheckRouter)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initCheckRouter)
	core.RegisterAdminObject(initAdminObject())
}
