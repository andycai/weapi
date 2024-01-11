package model

import (
	"gorm.io/gorm"
)

const TableNameComment = "comments"

type Comment struct {
	gorm.Model
	Body   string `gorm:"column:body;not null" json:"body"`
	Email  string `gorm:"column:email;not null" json:"email"`
	IP     string `gorm:"column:ip;not null" json:"ip"`
	Post   Post
	PostID uint
	UserID uint `gorm:"column:user_id;not null" json:"user_id"`
	User   User
}

func (*Comment) TableName() string {
	return TableNameComment
}
