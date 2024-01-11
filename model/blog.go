package model

import (
	"gorm.io/gorm"
)

const TableNameBlog = "blogs"

type Blog struct {
	gorm.Model
	Name        string `gorm:"column:name;not null" json:"name"`
	Description string `gorm:"column:description;not null" json:"description"`
	UserID      uint   `gorm:"column:user_id;not null" json:"user_id"`
	User        User
}

func (*Blog) TableName() string {
	return TableNameBlog
}
