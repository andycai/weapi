package model

import (
	"gorm.io/gorm"
)

const TableNameTag = "tags"

type Tag struct {
	gorm.Model
	Slug  string `gorm:"column:slug;not null" json:"slug"`
	Name  string `gorm:"column:name;not null;uniqueIndex" json:"name"`
	Posts []Post `gorm:"many2many:post_tag"`
}

func (*Tag) TableName() string {
	return TableNameTag
}
