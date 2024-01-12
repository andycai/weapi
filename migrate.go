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
			&model.PostTag{},
			&model.Page{},
			&model.Tag{},
			&model.Comment{},
			&model.Blog{},
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
		&model.PostTag{},
		&model.Page{},
		&model.Tag{},
		&model.Comment{},
		&model.Blog{},
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
