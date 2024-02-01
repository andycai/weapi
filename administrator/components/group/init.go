package group

import (
	"github.com/andycai/weapi"
	"github.com/andycai/weapi/administrator/components/user"
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.group.gorm.db"
	KeyNoCheckRouter = "admin.group.router.nocheck"
	KeyCheckRouter   = "admin.group.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initCheckRouter(r fiber.Router) {
}

func initAdminObject() []model.AdminObject {
	return []model.AdminObject{
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
			Attributes: map[string]model.AdminAttribute{
				"Role": {
					Default: model.GroupRoleMember,
					Choices: []model.AdminSelectOption{{Label: "Admin", Value: model.GroupRoleAdmin}, {Label: "Member", Value: model.GroupRoleMember}},
				},
			},
			Weight: 23,
		},
	}
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initCheckRouter)
	core.RegisterAdminObject(initAdminObject())
}
