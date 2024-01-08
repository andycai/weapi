package core

import (
	"github.com/andycai/werite/model"
	"gorm.io/gorm"
)

func AutoMigrate(dbs []*gorm.DB) {
	for _, db := range dbs {
		db.AutoMigrate(
			&model.User{},
			&model.Post{},
			&model.Category{},
			&model.PostTag{},
			&model.Page{},
			&model.Tag{},
			&model.Comment{},
			&model.Blog{},
		)
	}
}
