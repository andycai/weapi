package entity

import (
	"github.com/andycai/weapi"
	"github.com/andycai/weapi/administrator/components/user"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/andycai/weapi/object"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.object.gorm.db"
	KeyNoCheckRouter = "admin.object.router.nocheck"
	KeyCheckRouter   = "admin.object.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initNoCheckRouter(r fiber.Router) {
}

func initCheckRouter(r fiber.Router) {
	adminObjects := BuildAdminObjects(r, core.GetAdminObjects())

	r.Post("/json", func(c *fiber.Ctx) error {
		return JsonAction(c, adminObjects)
	})

	r.Post("/summary", HandleAdminSummary)
	r.Post("/tags/:content_type", handleGetTags)
}

func initAdminObject() []object.AdminObject {
	return []object.AdminObject{
		{
			Model:       &model.Group{},
			Group:       "Settings",
			Name:        "Group",
			Desc:        "A group describes a group of users. One user can be part of many groups and one group can have many users", //
			PluralName:  "Groups",
			Shows:       []string{"ID", "Name", "Extra", "UpdatedAt", "CreatedAt"},
			Editables:   []string{"ID", "Name", "Extra", "UpdatedAt"},
			Orderables:  []string{"UpdatedAt"},
			Searchables: []string{"Name"},
			Requireds:   []string{"Name"},
			Icon:        weapi.ReadIcon("/icon/group.svg"),
			AccessCheck: user.SuperAccessCheck,
			Weight:      22,
		},
		{
			Model:       &model.GroupMember{},
			Group:       "Settings",
			Name:        "GroupMember",
			Desc:        "Group members", //
			PluralName:  "GroupMembers",
			Shows:       []string{"ID", "User", "Group", "Role", "CreatedAt"},
			Filterables: []string{"Group", "Role", "CreatedAt"},
			Editables:   []string{"ID", "User", "Group", "Role"},
			Orderables:  []string{"CreatedAt"},
			Searchables: []string{"User", "Group"},
			Requireds:   []string{"User", "Group", "Role"},
			Icon:        weapi.ReadIcon("/icon/members.svg"),
			AccessCheck: user.SuperAccessCheck,
			Attributes: map[string]object.AdminAttribute{
				"Role": {
					Default: model.GroupRoleMember,
					Choices: []object.AdminSelectOption{{Label: "Admin", Value: model.GroupRoleAdmin}, {Label: "Member", Value: model.GroupRoleMember}},
				},
			},
			Weight: 23,
		},
		{
			Model:       &model.Config{},
			Group:       "Settings",
			Name:        "Config",
			Desc:        "System config with database backend, You can change it in admin page, and it will take effect immediately without restarting the server", //
			PluralName:  "Configs",
			Shows:       []string{"Key", "Value", "Desc"},
			Editables:   []string{"Key", "Value", "Desc"},
			Orderables:  []string{"Key"},
			Searchables: []string{"Key", "Value", "Desc"},
			Requireds:   []string{"Key", "Value"},
			Icon:        weapi.ReadIcon("/icon/config.svg"),
			AccessCheck: user.SuperAccessCheck,
			Weight:      24,
		},
	}
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterNoCheckRouter(KeyNoCheckRouter, initNoCheckRouter)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initCheckRouter)
	core.RegisterAdminObject(initAdminObject())
}
