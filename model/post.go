package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNamePost = "posts"

type Post struct {
	gorm.Model
	Slug        string    `gorm:"column:slug;not null;uniqueIndex" json:"slug"`
	Title       string    `gorm:"column:title;not null" json:"title"`
	Description string    `gorm:"column:description;not null" json:"description"`
	Body        string    `gorm:"column:body;not null" json:"body"`
	CategoryID  uint      `gorm:"column:category_id" json:"category_id"`
	IsDraft     uint      `gorm:"column:is_draft" json:"is_draft"`
	PublishedAt time.Time `gorm:"column:published_at" json:"published_at"`
	UserID      uint      `gorm:"column:user_id;not null" json:"user_id"`
	User        User
	Category    Category
	Comments    []Comment
	Tags        []Tag `gorm:"many2many:post_tag"`
}

func (*Post) TableName() string {
	return TableNamePost
}

func (p Post) GetTagsAsCommaSeparated() string {
	return p.GetTagsAsCharSeparated(",")
}

func (p Post) GetTagsAsCharSeparated(sep string) string {
	tagsText := ""

	for i := 0; i < len(p.Tags); i++ {
		tagsText += p.Tags[i].Name + sep
	}

	return tagsText
}
