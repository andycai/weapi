package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNamePage = "pages"

type Page struct {
	gorm.Model
	Slug        string    `gorm:"column:slug;not null" json:"slug"`
	Title       string    `gorm:"column:title;not null" json:"title"`
	Body        string    `gorm:"column:body;not null" json:"body"`
	PublishedAt time.Time `gorm:"column:published_at" json:"published_at"`
	UserID      uint      `gorm:"column:user_id;not null" json:"user_id"`
	User        User
}

func (*Page) TableName() string {
	return TableNamePage
}
