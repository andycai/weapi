package model

import (
	"gorm.io/gorm"
)

const TableNameMedia = "medias"

type Media struct {
	gorm.Model
	Name        string `gorm:"column:name" json:"name"`
	Description string `gorm:"column:description" json:"description"`
	Path        string `gorm:"column:path;not null" json:"path"`
}

func (*Media) TableName() string {
	return TableNameMedia
}
