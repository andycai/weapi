package activity

import (
	"github.com/andycai/weapi/core"
	"github.com/andycai/weapi/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.activity.gorm.db"
	KeyNoCheckRouter = "admin.activity.router.nocheck"
	KeyCheckRouter   = "admin.activity.router.check"
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
			Model:       &model.Activity{},
			Group:       "Activities",
			Name:        "Activity",
			Desc:        "The activity data of the website can only be in JSON format",
			Shows:       []string{"ID", "Site", "Creator", "Club", "Name", "BeginAt", "EndAt"},
			Editables:   []string{"Site", "Club", "Creator", "Name", "Description", "BeginAt", "EndAt", "Kind", "Type", "Quota", "Waiting", "Stage", "FeeType", "Ahead", "Address"},
			Filterables: []string{"Site", "Club", "UpdatedAt", "BeginAt", "EndAt"},
			Orderables:  []string{"UpdatedAt"},
			Searchables: []string{"Name", "Description"},
			Requireds:   []string{"Name", "Description", "Kind", "Type", "FeeType", "Quota", "BeginAt", "EndAt"},
			// Icon:        weapi.ReadIcon("/icon/piece.svg"),
			Scripts: []model.AdminScript{
				{Src: "/static/admin/js/cms_activity.js", Onload: true},
			},
			Orders: []model.Order{
				{
					Name: "UpdatedAt",
					Op:   model.OrderOpDesc,
				},
			},
			Attributes: map[string]model.AdminAttribute{
				"Kind": {
					Default: model.ActivityKindBasketball,
					Choices: []model.AdminSelectOption{
						{Label: "Basketball", Value: model.ActivityKindBasketball},
						{Label: "Football", Value: model.ActivityKindFootball},
						{Label: "Volleyball", Value: model.ActivityKindVolleyball},
						{Label: "Badminton", Value: model.ActivityKindBadminton},
						{Label: "Dinner", Value: model.ActivityKindDinner},
						{Label: "Other", Value: model.ActivityKindOther},
					},
				},
				"Type": {
					Default: model.ActivityTypePublic,
					Choices: []model.AdminSelectOption{
						{Label: "Private", Value: model.ActivityTypePrivate},
						{Label: "Public", Value: model.ActivityTypePublic},
						{Label: "Club", Value: model.ActivityTypeClub},
					},
				},
				"FeeType": {
					Default: model.ActivityFeeTypeFree,
					Choices: []model.AdminSelectOption{
						{Label: "Free", Value: model.ActivityFeeTypeFree},
						{Label: "Fixed", Value: model.ActivityFeeTypeFixed},
						{Label: "AA", Value: model.ActivityFeeTypeAA},
						{Label: "Male Fixed", Value: model.ActivityFeeTypeMaleFixedFemaleAA},
						{Label: "Female Fixed", Value: model.ActivityFeeTypeMaleAAFemaleFixed},
					},
				},
			},
			Actions: []model.AdminAction{},
			BeforeCreate: func(ctx *fiber.Ctx, vptr any) error {
				return nil
			},
			BeforeUpdate: func(ctx *fiber.Ctx, vptr any, vals map[string]any) error {
				return nil
			},
			Weight: 20,
		},
		{
			Model:       &model.Club{},
			Group:       "Activities",
			Name:        "Club",
			Desc:        "The club data of the website can only be in JSON format",
			Shows:       []string{"ID", "Site", "Creator", "Name", "Level", "CreatedAt"},
			Editables:   []string{"Name", "Site", "Creator", "Level", "Logo", "Description", "Notice", "Addr"},
			Filterables: []string{"Site", "UpdatedAt"},
			Orderables:  []string{"UpdatedAt"},
			Searchables: []string{"Name", "Description"},
			Requireds:   []string{"Name"},
			// Icon:        weapi.ReadIcon("/icon/piece.svg"),
			Scripts: []model.AdminScript{
				{Src: "/static/admin/js/cms_club.js", Onload: true},
			},
			Orders: []model.Order{
				{
					Name: "UpdatedAt",
					Op:   model.OrderOpDesc,
				},
			},
			Actions: []model.AdminAction{},
			BeforeCreate: func(ctx *fiber.Ctx, vptr any) error {
				return nil
			},
			BeforeUpdate: func(ctx *fiber.Ctx, vptr any, vals map[string]any) error {
				return nil
			},
			Weight: 21,
		},
		{
			Model:       &model.ClubMember{},
			Group:       "Activities",
			Name:        "ClubMember",
			Desc:        "The club member data of the website can only be in JSON format",
			Shows:       []string{"ID", "User", "Club", "Position", "DisplayName", "EnterAt"},
			Editables:   []string{"DisplayName", "User", "Club", "Position", "Scores", "EnterAt"},
			Filterables: []string{"Club", "EnterAt"},
			Orderables:  []string{"EnterAt"},
			Searchables: []string{"DisplayName"},
			Requireds:   []string{"DisplayName", "Position"},
			// Icon:        weapi.ReadIcon("/icon/piece.svg"),
			Attributes: map[string]model.AdminAttribute{
				"Position": {
					Default: model.ClubPositionMember,
					Choices: []model.AdminSelectOption{
						{Label: "Owner", Value: model.ClubPositionOwner},
						{Label: "Member", Value: model.ClubPositionMember},
						{Label: "Manager", Value: model.ClubPositionManager},
					},
				},
			},
			Actions: []model.AdminAction{},
			BeforeCreate: func(ctx *fiber.Ctx, vptr any) error {
				return nil
			},
			BeforeUpdate: func(ctx *fiber.Ctx, vptr any, vals map[string]any) error {
				return nil
			},
			Weight: 22,
		},
	}
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAdminObject(initAdminObject())
}
