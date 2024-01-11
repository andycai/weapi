package model

import (
	"gorm.io/gorm"
)

const TableNameCategory = "categories"

type Category struct {
	gorm.Model
	Slug        string `gorm:"column:slug;not null" json:"slug"`
	Name        string `gorm:"column:name;not null" json:"name"`
	Description string `gorm:"column:description;not null" json:"description"`
}

func (*Category) TableName() string {
	return TableNameCategory
}
