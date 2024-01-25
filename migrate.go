package weapi

import (
	"github.com/andycai/weapi/model"
	"gorm.io/gorm"
)

func AutoMultiMigrate(dbs []*gorm.DB) error {
	for _, db := range dbs {
		if err := MakeMigrates(db, []any{
			&model.User{},
			&model.Post{},
			&model.Category{},
			&model.Page{},
			&model.Comment{},
			&model.Config{},
			&model.Group{},
			&model.GroupMember{},
			&model.Site{},
			&model.Media{},
		}); err != nil {
			return err
		}
	}
	return nil
}

func AutoMigrate(db *gorm.DB) error {
	return MakeMigrates(db, []any{
		&model.User{},
		&model.Post{},
		&model.Category{},
		&model.Page{},
		&model.Comment{},
		&model.Config{},
		&model.Group{},
		&model.GroupMember{},
		&model.Site{},
		&model.Media{},
	})
}

func MakeMigrates(db *gorm.DB, insts []any) error {
	for _, v := range insts {
		if err := db.AutoMigrate(v); err != nil {
			return err
		}
	}
	return nil
}
