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
				return checkRequest(vptr.(*model.Activity))
			},
			BeforeUpdate: func(ctx *fiber.Ctx, vptr any, vals map[string]any) error {
				return checkRequest(vptr.(*model.Activity))
			},
			Weight: 21,
		},
	}
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAdminObject(initAdminObject())
}
