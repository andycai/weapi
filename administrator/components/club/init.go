package club

import (
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.club.gorm.db"
	KeyNoCheckRouter = "admin.club.router.nocheck"
	KeyCheckRouter   = "admin.club.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initAdminCheckRouter(r fiber.Router) {
}

func initAdminObject() []model.AdminObject {
	return []model.AdminObject{
		{
			Model:       &model.Club{},
			Group:       "Activities",
			Name:        "Club",
			Desc:        "The club data of the website can only be in JSON format",
			Shows:       []string{"ID", "Creator", "Name", "Level", "CreatedAt"},
			Editables:   []string{"Name", "Creator", "Level", "Logo", "Description", "Notice", "Addr"},
			Filterables: []string{"UpdatedAt"},
			Orderables:  []string{"UpdatedAt"},
			Searchables: []string{"Name", "Description"},
			Requireds:   []string{"Name"},
			// Icon:        weapi.ReadIcon("/icon/piece.svg"),
			Actions: []model.AdminAction{},
			BeforeCreate: func(ctx *fiber.Ctx, vptr any) error {
				return nil
			},
			BeforeUpdate: func(ctx *fiber.Ctx, vptr any, vals map[string]any) error {
				return nil
			},
			Weight: 22,
		},
		{
			Model:       &model.ClubMember{},
			Group:       "Activities",
			Name:        "ClubMember",
			Desc:        "The club member data of the website can only be in JSON format",
			Shows:       []string{"ID", "User", "Club", "Position", "DisplayName", "EnterAt"},
			Editables:   []string{"DisplayName", "User", "Club", "Position", "Scores"},
			Filterables: []string{"EnterAt"},
			Orderables:  []string{"EnterAt"},
			Searchables: []string{"DisplayName"},
			Requireds:   []string{"DisplayName"},
			// Icon:        weapi.ReadIcon("/icon/piece.svg"),
			Actions: []model.AdminAction{},
			BeforeCreate: func(ctx *fiber.Ctx, vptr any) error {
				return nil
			},
			BeforeUpdate: func(ctx *fiber.Ctx, vptr any, vals map[string]any) error {
				return nil
			},
			Weight: 23,
		},
	}
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAdminObject(initAdminObject())
}
